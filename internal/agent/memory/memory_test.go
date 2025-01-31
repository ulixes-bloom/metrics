package memory

import (
	"testing"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

var (
	gaugeValue   = float64(1)
	counterValue = int64(1)
)

func BenchmarkSetMetrics(b *testing.B) {
	for range b.N {
		s := NewStorage()
		for _, v := range metrics.GaugeMetrics {
			s.Set(metrics.NewGaugeMetric(v, gaugeValue))
		}
		for _, v := range metrics.CounterMetrics {
			s.Set(metrics.NewCounterMetric(v, counterValue))
		}
	}
}

func BenchmarkSetAll(b *testing.B) {
	metricsToSet := []metrics.Metric{}
	for _, v := range metrics.GaugeMetrics {
		metricsToSet = append(metricsToSet, metrics.NewGaugeMetric(v, gaugeValue))
	}
	for _, v := range metrics.CounterMetrics {
		metricsToSet = append(metricsToSet, metrics.NewCounterMetric(v, counterValue))
	}

	for range b.N {
		s := NewStorage()
		s.SetAll(metricsToSet)
	}
}
