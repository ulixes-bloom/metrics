package memory

import (
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type storage struct {
	metrics map[string]metrics.Metric
}

func NewStorage() *storage {
	m := storage{}
	m.metrics = map[string]metrics.Metric{}
	return &m
}

func (s *storage) Add(metric metrics.Metric) error {
	switch metric.MType {
	case metrics.Counter:
		cur, ok := s.metrics[metric.ID]
		if ok {
			newDelta := (metric.GetDelta() + cur.GetDelta())
			metric.Delta = &newDelta
			s.metrics[metric.ID] = metric
		} else {
			s.metrics[metric.ID] = metric
		}
	case metrics.Gauge:
		s.metrics[metric.ID] = metric
	default:
		return nil
	}

	return nil
}

func (s *storage) GetAll() map[string]metrics.Metric {
	return s.metrics
}
