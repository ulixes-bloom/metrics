package service

import (
	"math/rand"
	"runtime"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type service struct {
	storage Storage
}

func New(storage Storage) *service {
	return &service{storage: storage}
}

func (s *service) UpdateMetrics() {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	s.storage.Add(*metrics.NewGaugeMetric("Alloc", float64(ms.Alloc)))
	s.storage.Add(*metrics.NewGaugeMetric("BuckHashSys", float64(ms.BuckHashSys)))
	s.storage.Add(*metrics.NewGaugeMetric("Frees", float64(ms.Frees)))
	s.storage.Add(*metrics.NewGaugeMetric("GCCPUFraction", float64(ms.GCCPUFraction)))
	s.storage.Add(*metrics.NewGaugeMetric("GCSys", float64(ms.GCSys)))
	s.storage.Add(*metrics.NewGaugeMetric("HeapAlloc", float64(ms.HeapAlloc)))
	s.storage.Add(*metrics.NewGaugeMetric("HeapIdle", float64(ms.HeapIdle)))
	s.storage.Add(*metrics.NewGaugeMetric("HeapInuse", float64(ms.HeapInuse)))
	s.storage.Add(*metrics.NewGaugeMetric("HeapObjects", float64(ms.HeapObjects)))
	s.storage.Add(*metrics.NewGaugeMetric("HeapReleased", float64(ms.HeapReleased)))
	s.storage.Add(*metrics.NewGaugeMetric("HeapSys", float64(ms.HeapSys)))
	s.storage.Add(*metrics.NewGaugeMetric("LastGC", float64(ms.LastGC)))
	s.storage.Add(*metrics.NewGaugeMetric("Lookups", float64(ms.Lookups)))
	s.storage.Add(*metrics.NewGaugeMetric("MCacheInuse", float64(ms.MCacheInuse)))
	s.storage.Add(*metrics.NewGaugeMetric("MCacheSys", float64(ms.MCacheSys)))
	s.storage.Add(*metrics.NewGaugeMetric("MSpanInuse", float64(ms.MSpanInuse)))
	s.storage.Add(*metrics.NewGaugeMetric("MSpanSys", float64(ms.MSpanSys)))
	s.storage.Add(*metrics.NewGaugeMetric("Mallocs", float64(ms.Mallocs)))
	s.storage.Add(*metrics.NewGaugeMetric("NextGC", float64(ms.NextGC)))
	s.storage.Add(*metrics.NewGaugeMetric("NumForcedGC", float64(ms.NumForcedGC)))
	s.storage.Add(*metrics.NewGaugeMetric("NumGC", float64(ms.NumGC)))
	s.storage.Add(*metrics.NewGaugeMetric("OtherSys", float64(ms.OtherSys)))
	s.storage.Add(*metrics.NewGaugeMetric("PauseTotalNs", float64(ms.PauseTotalNs)))
	s.storage.Add(*metrics.NewGaugeMetric("StackInuse", float64(ms.StackInuse)))
	s.storage.Add(*metrics.NewGaugeMetric("StackSys", float64(ms.StackSys)))
	s.storage.Add(*metrics.NewGaugeMetric("Sys", float64(ms.Sys)))
	s.storage.Add(*metrics.NewGaugeMetric("TotalAlloc", float64(ms.TotalAlloc)))
	s.storage.Add(*metrics.NewGaugeMetric("RandomValue", rand.Float64()))

	s.storage.Add(*metrics.NewCounterMetric("PollCount", 1))
}

func (s *service) GetAll() map[string]metrics.Metric {
	return s.storage.GetAll()
}
