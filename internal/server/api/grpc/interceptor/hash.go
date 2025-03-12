package interceptor

import (
	"context"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/hash"
	"github.com/ulixes-bloom/ya-metrics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithHashing is a gRPC server-side interceptor that checks the hash of the incoming request's data.
// It compares the hash of the "Metric" in the request with the provided hash. If they do not match, the request is rejected.
func WithHashing(hashKey string) func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		updateMetricRequest, ok := req.(*proto.UpdateMetricRequest)
		if !ok {
			return nil, status.Error(codes.Internal, "Unable to parse updateMetricRequest")
		}

		h, err := hash.Encode([]byte(updateMetricRequest.Metric.String()), hashKey)
		if err != nil {
			return nil, status.Error(codes.Internal, "Unable to get hash for incomming metric")
		}

		if h != updateMetricRequest.GetHash() {
			return nil, status.Error(codes.PermissionDenied, "Incorrect hash")
		}

		return handler(ctx, req)
	}
}
