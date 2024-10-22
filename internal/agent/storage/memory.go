package storage

import (
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type MemStorage struct {
	Metrics map[string]metrics.Metric
}

func NewMemStorage() *MemStorage {
	m := MemStorage{}
	m.Metrics = make(map[string]metrics.Metric,
		len(metrics.GaugeMetrics)+len(metrics.CounterMetrics))

	return &m
}

func (m *MemStorage) Add(metric metrics.Metric) error {
	switch metric.MType {
	case metrics.Counter:
		cur, ok := m.Metrics[metric.ID]
		if ok {
			newDelta := (*metric.Delta + *cur.Delta)
			metric.Delta = &newDelta
			m.Metrics[metric.ID] = metric
		} else {
			m.Metrics[metric.ID] = metric
		}
	case metrics.Gauge:
		m.Metrics[metric.ID] = metric
	default:
		return nil
	}

	return nil
}

func (m *MemStorage) Get(name string) (val metrics.Metric, ok bool) {
	val, ok = m.Metrics[name]
	return
}

func (m *MemStorage) GetAll() map[string]metrics.Metric {
	return m.Metrics
}
