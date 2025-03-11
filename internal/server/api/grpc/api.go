package grpcapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api/grpc/interceptor"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
	"github.com/ulixes-bloom/ya-metrics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type grpcAPI struct {
	proto.UnimplementedMonitoringServer

	service api.Service
	conf    *config.Config
}

func (g *grpcAPI) UpdateMetric(ctx context.Context, in *proto.UpdateMetricRequest) (*proto.EmprtyResponse, error) {
	var EmprtyResponse proto.EmprtyResponse
	metric, err := metrics.NewMetric(in.Metric.Id, in.Metric.Mtype, in.Metric.GetValue(), in.Metric.GetDelta())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if _, err = g.service.UpdateJSONMetric(ctx, metric); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &EmprtyResponse, nil
}

func New(conf *config.Config, storage service.Storage) *grpcAPI {
	srv := service.New(storage, conf)
	newAPI := grpcAPI{
		service: srv,
		conf:    conf,
	}
	return &newAPI
}

func (g *grpcAPI) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	listen, err := net.Listen("tcp", g.conf.GRPCRunAddr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption

	// Chain interceptors
	var interceptors []grpc.UnaryServerInterceptor
	interceptors = append(interceptors, interceptor.WithLogging)

	// Add IP Resolving interceptor if the Trusted Subnet is set
	if g.conf.TrustedSubnet != "" {
		interceptors = append(interceptors, interceptor.WithIPResolving(g.conf.TrustedSubnet))
	}

	// Add Hashing interceptor if the HashKey is set
	if g.conf.HashKey != "" {
		interceptors = append(interceptors, interceptor.WithHashing(g.conf.HashKey))
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(interceptors...))

	// configure TLS
	if g.conf.PublicKey != "" && g.conf.PrivateKey != "" {
		creds, err := g.loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("grpcapi.run: %w", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	// Create new grpc server and register Monitoring Server implementation
	s := grpc.NewServer(opts...)
	proto.RegisterMonitoringServer(s, g)

	go func() {
		errChan <- s.Serve(listen)
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("grpcapi.run: %w", err)
	case <-ctx.Done():
		return g.service.Shutdown(ctx)
	}
}

// Load TLS credentials from PEM files
func (g *grpcAPI) loadTLSCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(g.conf.PublicKey, g.conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("grpcapi.loadTLSCredentials: failed to load key pair: %w", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	return creds, nil
}
