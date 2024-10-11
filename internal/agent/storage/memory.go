package storage

import (
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMemStorage() *MemStorage {
	m := MemStorage{}
	m.Gauges = make(map[string]float64, len(metrics.GaugeMetrics))
	m.Counters = make(map[string]int64, len(metrics.CounterMetrics))
	return &m
}

func (m *MemStorage) Add(name string, value interface{}) error {
	switch v := value.(type) {
	case int64:
		m.Counters[name] += v
	case float64:
		m.Gauges[name] = v
	default:
		return nil
	}

	return nil
}

func (m *MemStorage) GetGauge(name string) (val float64, ok bool) {
	val, ok = m.Gauges[name]
	return
}

func (m *MemStorage) GetCounter(name string) (val int64, ok bool) {
	val, ok = m.Counters[name]
	return
}

func (m *MemStorage) GetAllGauges() map[string]float64 {
	return m.Gauges
}

func (m *MemStorage) GetAllCounters() map[string]int64 {
	return m.Counters
}
