// Package grpcclient provides GRPC implemetation of the agent
// It is ressponsible for polling system metrics and
// reporting them to a remote server via GRPC.
package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/client"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/hash"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/workerpool"
	"github.com/ulixes-bloom/ya-metrics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// client handles polling metrics from the system and reporting them to a server.
type grpcClient struct {
	service client.Service
	conn    *grpc.ClientConn
	conf    *config.Config
	ip      string
}

// New creates and initializes a new client instance.
func New(conf *config.Config, storage service.Storage) (*grpcClient, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("grpcclient.new: Error while retrieving interface addresses, %w", err)
	}
	var ip string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	// generate credentials for grpc connection
	creds := insecure.NewCredentials()
	if conf.CryptoKey != "" {
		pemData, err := os.ReadFile(conf.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("grpcclient.loadTLSCredentials: Failed to read public.pem: %v", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemData) {
			return nil, fmt.Errorf("grpcclient.loadTLSCredentials: Failed to add public key to cert pool")
		}

		creds = credentials.NewTLS(&tls.Config{
			RootCAs: certPool,
		})
	}

	// open grpc connection with server
	conn, err := grpc.NewClient(conf.ServerAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("grpcclient.new: Error while creating grpc connection, %w", err)
	}

	return &grpcClient{
		service: service.New(storage),
		conn:    conn,
		conf:    conf,
		ip:      ip,
	}, nil
}

// Run starts background operations of the client:
//
// 1. Polling metrics with the period specified in config.PollInterval.
//
// 2. Reporting metrics to the server with the period specified in config.ReportInterval.
//
// It runs these operations concurrently and waits for them to complete.
func (c *grpcClient) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		c.pollMetrics(ctx)
		wg.Done()
	}()

	go func() {
		c.reportMetrics(ctx)
		wg.Done()
	}()

	wg.Wait()
}

// pollMetrics periodically polls system metrics and stores them in memory.
func (c *grpcClient) pollMetrics(ctx context.Context) {
	pollTicker := time.NewTicker(c.conf.GetPollIntervalDuration())
	defer pollTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			err := c.service.Poll(ctx)
			if err != nil {
				log.Error().Msg(err.Error())
			}
		case <-ctx.Done():
			log.Debug().Msg("done polling metrics")
			return
		}
	}
}

// reportMetrics periodically retrieves all stored metrics from memory and sends them to the server.
// sending is done in parallel using a workerpool.
func (c *grpcClient) reportMetrics(ctx context.Context) {
	reportTicker := time.NewTicker(c.conf.GetReportIntervalDuration())
	defer reportTicker.Stop()

	// create worker pool for sending metrics to server
	pool := workerpool.New(c.conf.RateLimit, metrics.MetricsCount, c.sendMetric)

	for {
		select {
		case <-reportTicker.C:
			for _, m := range c.service.GetAll() {
				pool.Submit(m)
			}
		case <-ctx.Done():
			log.Debug().Msg("done reporting metrics")
			pool.StopAndWait()
			c.conn.Close()
			return
		}
	}
}

// sendMetric sends a single metric to the server after compressing and encoding it.
func (c *grpcClient) sendMetric(m metrics.Metric) error {
	client := proto.NewMonitoringClient(c.conn)
	var updateMetricRequest proto.UpdateMetricRequest
	pbMetric := &proto.Metric{
		Id:    m.ID,
		Mtype: m.MType,
		Value: m.Value,
		Delta: m.Delta,
	}
	updateMetricRequest.Metric = pbMetric

	// set agent ip in grpc request metadata
	md := metadata.New(map[string]string{"x-real-ip": c.ip})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// calculate and set metric hash in request
	if c.conf.HashKey != "" {
		h, err := hash.Encode([]byte(pbMetric.String()), c.conf.HashKey)
		if err != nil {
			return fmt.Errorf("grpcclient.sendMetric: %w", err)
		}
		updateMetricRequest.Hash = &h
	}

	_, err := client.UpdateMetric(ctx, &updateMetricRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return fmt.Errorf("grpcclient.sendMetric: %s, %s", e.Message(), e.Code())
		} else {
			return fmt.Errorf("grpcclient.sendMetric: Can't parse error %w", err)
		}
	}

	return nil
}

func (c *grpcClient) loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemData, err := os.ReadFile(c.conf.CryptoKey)
	if err != nil {
		return nil, fmt.Errorf("grpcclient.loadTLSCredentials: Failed to read public.pem: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemData) {
		return nil, fmt.Errorf("grpcclient.loadTLSCredentials: Failed to add public key to cert pool")
	}

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})
	return creds, nil
}
