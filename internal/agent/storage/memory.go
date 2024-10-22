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

	for _, g := range metrics.GaugeMetrics {
		zeroVal := float64(0)
		m.Metrics[g] = metrics.Metric{
			ID:    g,
			MType: metrics.Gauge,
			Value: &zeroVal,
		}
	}
	for _, c := range metrics.CounterMetrics {
		zeroVal := int64(0)
		m.Metrics[c] = metrics.Metric{
			ID:    c,
			MType: metrics.Counter,
			Delta: &zeroVal,
		}
	}
	return &m
}

func (m *MemStorage) Add(metric metrics.Metric) error {
	switch metric.MType {
	case metrics.Counter:
		cur := m.Metrics[metric.ID]
		newDelta := (*metric.Delta + *cur.Delta)
		metric.Delta = &newDelta
		m.Metrics[metric.ID] = metric
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
