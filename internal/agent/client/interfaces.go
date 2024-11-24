package client

import (
	"context"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type Service interface {
	GetAll() map[string]metrics.Metric
	Poll(context.Context) error
}
