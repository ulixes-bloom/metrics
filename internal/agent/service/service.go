package service

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"golang.org/x/sync/errgroup"
)

type service struct {
	storage Storage
}

func New(storage Storage) *service {
	return &service{storage: storage}
}

func (s *service) Poll(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := s.pollRuntimeMetrics()
		if err != nil {
			return fmt.Errorf("pollRuntimeMetrics: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		err := s.pollSystemMetrics(ctx)
		if err != nil {
			return fmt.Errorf("pollSystemMetrics: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to poll metrics: %w", err)
	}
	return nil
}

func (s *service) pollRuntimeMetrics() error {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	metricsValues := []metrics.Metric{
		metrics.NewGaugeMetric("Alloc", float64(ms.Alloc)),
		metrics.NewGaugeMetric("BuckHashSys", float64(ms.BuckHashSys)),
		metrics.NewGaugeMetric("Frees", float64(ms.Frees)),
		metrics.NewGaugeMetric("GCCPUFraction", float64(ms.GCCPUFraction)),
		metrics.NewGaugeMetric("GCSys", float64(ms.GCSys)),
		metrics.NewGaugeMetric("HeapAlloc", float64(ms.HeapAlloc)),
		metrics.NewGaugeMetric("HeapIdle", float64(ms.HeapIdle)),
		metrics.NewGaugeMetric("HeapInuse", float64(ms.HeapInuse)),
		metrics.NewGaugeMetric("HeapObjects", float64(ms.HeapObjects)),
		metrics.NewGaugeMetric("HeapReleased", float64(ms.HeapReleased)),
		metrics.NewGaugeMetric("HeapSys", float64(ms.HeapSys)),
		metrics.NewGaugeMetric("LastGC", float64(ms.LastGC)),
		metrics.NewGaugeMetric("Lookups", float64(ms.Lookups)),
		metrics.NewGaugeMetric("MCacheInuse", float64(ms.MCacheInuse)),
		metrics.NewGaugeMetric("MCacheSys", float64(ms.MCacheSys)),
		metrics.NewGaugeMetric("MSpanInuse", float64(ms.MSpanInuse)),
		metrics.NewGaugeMetric("MSpanSys", float64(ms.MSpanSys)),
		metrics.NewGaugeMetric("Mallocs", float64(ms.Mallocs)),
		metrics.NewGaugeMetric("NextGC", float64(ms.NextGC)),
		metrics.NewGaugeMetric("NumForcedGC", float64(ms.NumForcedGC)),
		metrics.NewGaugeMetric("NumGC", float64(ms.NumGC)),
		metrics.NewGaugeMetric("OtherSys", float64(ms.OtherSys)),
		metrics.NewGaugeMetric("PauseTotalNs", float64(ms.PauseTotalNs)),
		metrics.NewGaugeMetric("StackInuse", float64(ms.StackInuse)),
		metrics.NewGaugeMetric("StackSys", float64(ms.StackSys)),
		metrics.NewGaugeMetric("Sys", float64(ms.Sys)),
		metrics.NewGaugeMetric("TotalAlloc", float64(ms.TotalAlloc)),
		metrics.NewGaugeMetric("RandomValue", rand.Float64()),

		metrics.NewCounterMetric("PollCount", 1),
	}

	err := s.storage.SetAll(metricsValues)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) pollSystemMetrics(ctx context.Context) error {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error while polling virtual memory: %w", err)
	}

	utilisation, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return fmt.Errorf("error while polling cpu: %w", err)
	}

	metricsValues := []metrics.Metric{
		metrics.NewGaugeMetric("TotalMemory", float64(vMem.Total)),
		metrics.NewGaugeMetric("FreeMemory", float64(vMem.Free)),
		metrics.NewGaugeMetric("CPUutilization1", float64(utilisation[0])),
	}

	err = s.storage.SetAll(metricsValues)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAll() map[string]metrics.Metric {
	return s.storage.GetAll()
}
