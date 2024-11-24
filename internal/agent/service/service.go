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
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.pollRuntimeMetrics()
		return nil
	})
	g.Go(func() error {
		return s.pollSystemMetrics(ctx)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to poll metrics: %w", err)
	}
	return nil
}

func (s *service) pollRuntimeMetrics() {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	s.storage.Add(metrics.NewGaugeMetric("Alloc", float64(ms.Alloc)))
	s.storage.Add(metrics.NewGaugeMetric("BuckHashSys", float64(ms.BuckHashSys)))
	s.storage.Add(metrics.NewGaugeMetric("Frees", float64(ms.Frees)))
	s.storage.Add(metrics.NewGaugeMetric("GCCPUFraction", float64(ms.GCCPUFraction)))
	s.storage.Add(metrics.NewGaugeMetric("GCSys", float64(ms.GCSys)))
	s.storage.Add(metrics.NewGaugeMetric("HeapAlloc", float64(ms.HeapAlloc)))
	s.storage.Add(metrics.NewGaugeMetric("HeapIdle", float64(ms.HeapIdle)))
	s.storage.Add(metrics.NewGaugeMetric("HeapInuse", float64(ms.HeapInuse)))
	s.storage.Add(metrics.NewGaugeMetric("HeapObjects", float64(ms.HeapObjects)))
	s.storage.Add(metrics.NewGaugeMetric("HeapReleased", float64(ms.HeapReleased)))
	s.storage.Add(metrics.NewGaugeMetric("HeapSys", float64(ms.HeapSys)))
	s.storage.Add(metrics.NewGaugeMetric("LastGC", float64(ms.LastGC)))
	s.storage.Add(metrics.NewGaugeMetric("Lookups", float64(ms.Lookups)))
	s.storage.Add(metrics.NewGaugeMetric("MCacheInuse", float64(ms.MCacheInuse)))
	s.storage.Add(metrics.NewGaugeMetric("MCacheSys", float64(ms.MCacheSys)))
	s.storage.Add(metrics.NewGaugeMetric("MSpanInuse", float64(ms.MSpanInuse)))
	s.storage.Add(metrics.NewGaugeMetric("MSpanSys", float64(ms.MSpanSys)))
	s.storage.Add(metrics.NewGaugeMetric("Mallocs", float64(ms.Mallocs)))
	s.storage.Add(metrics.NewGaugeMetric("NextGC", float64(ms.NextGC)))
	s.storage.Add(metrics.NewGaugeMetric("NumForcedGC", float64(ms.NumForcedGC)))
	s.storage.Add(metrics.NewGaugeMetric("NumGC", float64(ms.NumGC)))
	s.storage.Add(metrics.NewGaugeMetric("OtherSys", float64(ms.OtherSys)))
	s.storage.Add(metrics.NewGaugeMetric("PauseTotalNs", float64(ms.PauseTotalNs)))
	s.storage.Add(metrics.NewGaugeMetric("StackInuse", float64(ms.StackInuse)))
	s.storage.Add(metrics.NewGaugeMetric("StackSys", float64(ms.StackSys)))
	s.storage.Add(metrics.NewGaugeMetric("Sys", float64(ms.Sys)))
	s.storage.Add(metrics.NewGaugeMetric("TotalAlloc", float64(ms.TotalAlloc)))
	s.storage.Add(metrics.NewGaugeMetric("RandomValue", rand.Float64()))

	s.storage.Add(metrics.NewCounterMetric("PollCount", 1))
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

	s.storage.Add(metrics.NewGaugeMetric("TotalMemory", float64(vMem.Total)))
	s.storage.Add(metrics.NewGaugeMetric("FreeMemory", float64(vMem.Free)))
	s.storage.Add(metrics.NewGaugeMetric("CPUutilization1", float64(utilisation[0])))

	return nil
}

func (s *service) GetAll() map[string]metrics.Metric {
	return s.storage.GetAll()
}
