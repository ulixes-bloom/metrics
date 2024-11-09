package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Storage interface {
	Add(value metrics.Metric) error
	GetAll() map[string]metrics.Metric
}
