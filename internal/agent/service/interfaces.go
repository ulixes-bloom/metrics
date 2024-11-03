package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Storage interface {
	Add(value metrics.Metric) error
	Get(name string) (val metrics.Metric, ok bool)
	GetAll() map[string]metrics.Metric
}
