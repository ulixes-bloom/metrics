package service

import (
	"encoding/json"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type Service struct {
	Storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{Storage: storage}
}

func (s *Service) GetMetricsHTMLTable() ([]byte, error) {
	return s.Storage.HTMLTable()
}

func (s *Service) GetJSONMetric(mtype, mname string) ([]byte, error) {
	val, ok := s.Storage.Get(mname)
	if !ok {
		return []byte(""), errors.ErrMetricNotExists
	}
	return json.Marshal(val)
}

func (s *Service) UpdateJSONMetric(metric metrics.Metric) ([]byte, error) {
	metric, err := s.Storage.Add(metric)
	if err != nil {
		return []byte(""), err
	}
	return json.Marshal(metric)
}
