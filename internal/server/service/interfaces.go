package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Storage interface {
	Add(metric metrics.Metric) (metrics.Metric, error)
	Get(name string) (val metrics.Metric, ok bool)
	All() map[string]string
	HTMLTable() ([]byte, error)
}
