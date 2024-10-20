package service

import (
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
		val, ok := s.Storage.GetGauge(mname)
		if !ok {
			return []byte(""), errors.ErrMetricNotExists
		}
		mval = strconv.FormatFloat(val, 'f', -1, 64)
	case metrics.Counter:
		val, ok := s.Storage.GetCounter(mname)
		if !ok {
			return []byte(""), errors.ErrMetricNotExists
		}
		mval = strconv.FormatInt(val, 10)
	default:
		return []byte(""), errors.ErrMetricTypeNotImplemented
	}

	return []byte(mval), nil
}

func (s *Service) UpdateMetric(mtype, mname, mval string) error {
	switch mtype {
	case metrics.Gauge:
		if val, err := strconv.ParseFloat(mval, 64); err == nil {
			s.Storage.AddGauge(mname, val)
		} else {
			return errors.ErrMetricValueNotValid
		}
	case metrics.Counter:
		if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
			s.Storage.AddCounter(mname, val)
		} else {
			return errors.ErrMetricValueNotValid
		}
	default:
		return errors.ErrMetricTypeNotImplemented
	}

	return nil
}
