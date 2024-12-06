package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Storage interface {
	Set(value metrics.Metric) error
	SetAll(meticsSlice []metrics.Metric) error
	GetAll() map[string]metrics.Metric
}
