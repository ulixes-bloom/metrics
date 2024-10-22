package service

import (
	"encoding/json"
	"strconv"

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

func (s *Service) GetMetric(mtype, mname string) ([]byte, error) {
	var mval string

	switch mtype {
	case metrics.Gauge:
		metric, ok := s.Storage.Get(mname)
		if !ok {
			return []byte(""), errors.ErrMetricNotExists
		}
		mval = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case metrics.Counter:
		metric, ok := s.Storage.Get(mname)
		if !ok {
			return []byte(""), errors.ErrMetricNotExists
		}
		mval = strconv.FormatInt(*metric.Delta, 10)
	default:
		return []byte(""), errors.ErrMetricTypeNotImplemented
	}

	return []byte(mval), nil
}

func (s *Service) UpdateMetric(mtype, mname, mval string) error {
	switch mtype {
	case metrics.Gauge:
		if val, err := strconv.ParseFloat(mval, 64); err == nil {
			s.Storage.Add(*metrics.NewGaugeMetric(mname, val))
		} else {
			return errors.ErrMetricValueNotValid
		}
	case metrics.Counter:
		if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
			s.Storage.Add(*metrics.NewCounterMetric(mname, val))
		} else {
			return errors.ErrMetricValueNotValid
		}
	default:
		return errors.ErrMetricTypeNotImplemented
	}

	return nil
}

func (s *Service) GetJSONMetric(metric metrics.Metric) ([]byte, error) {
	val, ok := s.Storage.Get(metric.ID)
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
