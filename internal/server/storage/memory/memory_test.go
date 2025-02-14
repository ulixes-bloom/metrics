package memory

import (
	"context"
	"testing"
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
)

var (
	gaugeValue     = float64(1)
	counterValue   = int64(1)
	contextTimeout = 30 * time.Second
)

func BenchmarkSet(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	config := config.GetDefault()
	for range b.N {
		s, _ := NewStorage(ctx, config)
		for _, v := range metrics.GaugeMetrics {
			s.Set(ctx, metrics.NewGaugeMetric(v, gaugeValue))
		}
		for _, v := range metrics.CounterMetrics {
			s.Set(ctx, metrics.NewCounterMetric(v, counterValue))
		}
	}
}

func BenchmarkSetAll(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	config := config.GetDefault()

	metricsToSet := []metrics.Metric{}
	for _, v := range metrics.GaugeMetrics {
		metricsToSet = append(metricsToSet, metrics.NewGaugeMetric(v, gaugeValue))
	}
	for _, v := range metrics.CounterMetrics {
		metricsToSet = append(metricsToSet, metrics.NewCounterMetric(v, counterValue))
	}

	for range b.N {
		s, _ := NewStorage(ctx, config)
		s.SetAll(ctx, metricsToSet)
	}
}
