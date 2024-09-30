package storage

import (
	"github.com/ulixes-bloom/ya-metrics/internal/metrics"
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

func (m *MemStorage) AddGauge(name string, value float64) {
	m.Gauges[name] = value
}

func (m *MemStorage) AddCounter(name string, value int64) {
	m.Counters[name] += value
}
