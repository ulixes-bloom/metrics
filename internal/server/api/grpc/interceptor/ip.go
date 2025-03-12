package interceptor

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// WithIPResolving is a gRPC server-side interceptor for checking if the client's IP address is within a trusted subnet.
func WithIPResolving(trustedSubnet string) func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ipReq, err := resolveIP(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		_, ipnet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if !ipnet.Contains(ipReq) {
			return nil, status.Error(codes.PermissionDenied, "mUntrusted IP address")
		}

		return handler(ctx, req)
	}
}

// resolveIP extracts the client's IP address from the gRPC request context by reading the "x-real-ip" metadata.
// It returns the resolved IP address or an error if the IP is not found or cannot be parsed.
func resolveIP(ctx context.Context) (net.IP, error) {
	var ipStr string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("x-real-ip")
		if len(values) > 0 {
			ipStr = values[0]
		}
	}
	ip := net.ParseIP(ipStr)

	if ip == nil {
		return nil, fmt.Errorf("ip.resolveIP: failed to parse ip from http header")
	}
	return ip, nil
}
