package client

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Service interface {
	UpdateMetrics()
	GetAll() map[string]metrics.Metric
}
