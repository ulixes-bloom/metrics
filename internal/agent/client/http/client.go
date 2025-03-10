// Package httpclient provides HTTP implemetation of the agent
// It is ressponsible for polling system metrics and
// reporting them to a remote server via HTTP.
package httpclient

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/client"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/rsa"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/workerpool"
)

// httpClient handles polling metrics from the system and reporting them to a server.
type httpClient struct {
	service client.Service
	http    *http.Client
	conf    *config.Config
	ip      string
}

// New creates and initializes a new client instance.
func New(conf *config.Config, storage service.Storage) (*httpClient, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("client.new: Error while retrieving interface addresses, %w", err)
	}
	var ip string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return &httpClient{
		service: service.New(storage),
		http:    &http.Client{},
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
func (c *httpClient) Run(ctx context.Context) {
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
func (c *httpClient) pollMetrics(ctx context.Context) {
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
func (c *httpClient) reportMetrics(ctx context.Context) {
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
			return
		}
	}
}

// sendMetric sends a single metric to the server after compressing and encoding it.
func (c *httpClient) sendMetric(m metrics.Metric) error {
	// Marshal metric to JSON
	marshalled, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("client.sendMetric: %w", err)
	}

	if c.conf.CryptoKey != "" {
		marshalled, err = rsa.Encrypt(marshalled, c.conf.CryptoKey)
		if err != nil {
			return fmt.Errorf("client.sendMetric: %w", err)
		}
	}

	// Compress the marshalled data
	var buf bytes.Buffer
	gb := gzip.NewWriter(&buf)
	_, err = gb.Write(marshalled)
	if err != nil {
		return fmt.Errorf("client.sendMetric: failed to compress metric, %w", err)
	}
	err = gb.Close()
	if err != nil {
		return fmt.Errorf("client.sendMetric: failed to close gzip writer, %w", err)
	}

	// Construct the request
	url := fmt.Sprintf("%s/update/", c.conf.GetNormilizedServerAddr())
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return err
	}
	req.Header.Add(headers.ContentType, "application/json")
	req.Header.Add(headers.AcceptEncoding, "gzip")
	req.Header.Add(headers.ContentEncoding, "gzip")
	req.Header.Add(headers.XRealIP, c.ip)

	// Execute the HTTP request
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("client.sendMetric: %w", err)
	}
	defer res.Body.Close() // Ensure the body is always closed

	// Validate the response status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("client.sendMetric: unexpected response status code while sending metrics '%s'", res.Status)
	}

	return nil
}
