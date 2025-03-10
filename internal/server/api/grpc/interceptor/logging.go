package interceptor

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// WithLogging is a gRPC Unary Server interceptor for logging the details of incoming requests.
// It logs the method name and the time it takes to process the request.
func WithLogging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	method := info.FullMethod

	res, err := handler(ctx, req)

	duration := time.Since(start)
	log.Debug().
		Str("method", method).
		Str("duration", duration.String()).
		Msg("got incoming grpc request")

	return res, err
}
