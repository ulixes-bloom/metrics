package service

import (
	"math/rand"
	"runtime"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type Service struct {
	Storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{Storage: storage}
}

func (s *Service) UpdateMetrics() {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	s.Storage.Add(*metrics.NewGaugeMetric("Alloc", float64(ms.Alloc)))
	s.Storage.Add(*metrics.NewGaugeMetric("BuckHashSys", float64(ms.BuckHashSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("Frees", float64(ms.Frees)))
	s.Storage.Add(*metrics.NewGaugeMetric("GCCPUFraction", float64(ms.GCCPUFraction)))
	s.Storage.Add(*metrics.NewGaugeMetric("GCSys", float64(ms.GCSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("HeapAlloc", float64(ms.HeapAlloc)))
	s.Storage.Add(*metrics.NewGaugeMetric("HeapIdle", float64(ms.HeapIdle)))
	s.Storage.Add(*metrics.NewGaugeMetric("HeapInuse", float64(ms.HeapInuse)))
	s.Storage.Add(*metrics.NewGaugeMetric("HeapObjects", float64(ms.HeapObjects)))
	s.Storage.Add(*metrics.NewGaugeMetric("HeapReleased", float64(ms.HeapReleased)))
	s.Storage.Add(*metrics.NewGaugeMetric("HeapSys", float64(ms.HeapSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("LastGC", float64(ms.LastGC)))
	s.Storage.Add(*metrics.NewGaugeMetric("Lookups", float64(ms.Lookups)))
	s.Storage.Add(*metrics.NewGaugeMetric("MCacheInuse", float64(ms.MCacheInuse)))
	s.Storage.Add(*metrics.NewGaugeMetric("MCacheSys", float64(ms.MCacheSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("MSpanInuse", float64(ms.MSpanInuse)))
	s.Storage.Add(*metrics.NewGaugeMetric("MSpanSys", float64(ms.MSpanSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("Mallocs", float64(ms.Mallocs)))
	s.Storage.Add(*metrics.NewGaugeMetric("NextGC", float64(ms.NextGC)))
	s.Storage.Add(*metrics.NewGaugeMetric("NumForcedGC", float64(ms.NumForcedGC)))
	s.Storage.Add(*metrics.NewGaugeMetric("NumGC", float64(ms.NumGC)))
	s.Storage.Add(*metrics.NewGaugeMetric("OtherSys", float64(ms.OtherSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("PauseTotalNs", float64(ms.PauseTotalNs)))
	s.Storage.Add(*metrics.NewGaugeMetric("StackInuse", float64(ms.StackInuse)))
	s.Storage.Add(*metrics.NewGaugeMetric("StackSys", float64(ms.StackSys)))
	s.Storage.Add(*metrics.NewGaugeMetric("Sys", float64(ms.Sys)))
	s.Storage.Add(*metrics.NewGaugeMetric("TotalAlloc", float64(ms.TotalAlloc)))
	s.Storage.Add(*metrics.NewGaugeMetric("RandomValue", rand.Float64()))

	s.Storage.Add(*metrics.NewCounterMetric("PollCount", 1))
}

func (s *Service) GetAll() map[string]metrics.Metric {
	return s.Storage.GetAll()
}
