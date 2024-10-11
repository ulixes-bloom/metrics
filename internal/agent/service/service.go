package service

import (
	"math/rand"
	"runtime"
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

	s.Storage.Add("Alloc", float64(ms.Alloc))
	s.Storage.Add("BuckHashSys", float64(ms.BuckHashSys))
	s.Storage.Add("Frees", float64(ms.Frees))
	s.Storage.Add("GCCPUFraction", float64(ms.GCCPUFraction))
	s.Storage.Add("GCSys", float64(ms.GCSys))
	s.Storage.Add("HeapAlloc", float64(ms.HeapAlloc))
	s.Storage.Add("HeapIdle", float64(ms.HeapIdle))
	s.Storage.Add("HeapInuse", float64(ms.HeapInuse))
	s.Storage.Add("HeapObjects", float64(ms.HeapObjects))
	s.Storage.Add("HeapReleased", float64(ms.HeapReleased))
	s.Storage.Add("HeapSys", float64(ms.HeapSys))
	s.Storage.Add("LastGC", float64(ms.LastGC))
	s.Storage.Add("Lookups", float64(ms.Lookups))
	s.Storage.Add("MCacheInuse", float64(ms.MCacheInuse))
	s.Storage.Add("MCacheSys", float64(ms.MCacheSys))
	s.Storage.Add("MSpanInuse", float64(ms.MSpanInuse))
	s.Storage.Add("MSpanSys", float64(ms.MSpanSys))
	s.Storage.Add("Mallocs", float64(ms.Mallocs))
	s.Storage.Add("NextGC", float64(ms.NextGC))
	s.Storage.Add("NumForcedGC", float64(ms.NumForcedGC))
	s.Storage.Add("NumGC", float64(ms.NumGC))
	s.Storage.Add("OtherSys", float64(ms.OtherSys))
	s.Storage.Add("PauseTotalNs", float64(ms.PauseTotalNs))
	s.Storage.Add("StackInuse", float64(ms.StackInuse))
	s.Storage.Add("StackSys", float64(ms.StackSys))
	s.Storage.Add("Sys", float64(ms.Sys))
	s.Storage.Add("TotalAlloc", float64(ms.TotalAlloc))
	s.Storage.Add("RandomValue", rand.Float64())

	s.Storage.Add("PollCount", 1)
}

func (s *Service) GetAllCounters() map[string]int64 {
	return s.Storage.GetAllCounters()
}

func (s *Service) GetAllGauges() map[string]float64 {
	return s.Storage.GetAllGauges()
}
