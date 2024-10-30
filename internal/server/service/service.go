package service

import (
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
)

type service struct {
	storage Storage
	log     zerolog.Logger
}

func New(storage Storage, logger zerolog.Logger, conf config.Config) *service {
	srv := &service{storage: storage, log: logger}

	if conf.Restore {
		err := srv.RestoreMetrics()
		if err != nil {
			srv.log.Error().Msg(err.Error())
		}
	}
	return srv
}

func (s *service) GetMetricsHTMLTable() ([]byte, error) {
	return s.storage.HTMLTable()
}

func (s *service) GetMetric(mtype, mname string) ([]byte, error) {
	var mval string

	switch mtype {
	case metrics.Gauge:
		metric, ok := s.storage.Get(mname)
		if !ok {
			return []byte(""), errors.ErrMetricNotExists
		}
		mval = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case metrics.Counter:
		metric, ok := s.storage.Get(mname)
		if !ok {
			return []byte(""), errors.ErrMetricNotExists
		}
		mval = strconv.FormatInt(*metric.Delta, 10)
	default:
		return []byte(""), errors.ErrMetricTypeNotImplemented
	}

	return []byte(mval), nil
}

func (s *service) UpdateMetric(mtype, mname, mval string) error {
	switch mtype {
	case metrics.Gauge:
		if val, err := strconv.ParseFloat(mval, 64); err == nil {
			s.storage.Add(*metrics.NewGaugeMetric(mname, val))
		} else {
			return errors.ErrMetricValueNotValid
		}
	case metrics.Counter:
		if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
			s.storage.Add(*metrics.NewCounterMetric(mname, val))
		} else {
			return errors.ErrMetricValueNotValid
		}
	default:
		return errors.ErrMetricTypeNotImplemented
	}

	return nil
}

func (s *service) GetJSONMetric(metric metrics.Metric) ([]byte, error) {
	val, ok := s.storage.Get(metric.ID)
	if !ok {
		return []byte(""), errors.ErrMetricNotExists
	}

	return json.Marshal(val)
}

func (s *service) UpdateJSONMetric(metric metrics.Metric) ([]byte, error) {
	metric, err := s.storage.Add(metric)
	if err != nil {
		return []byte(""), err
	}
	return json.Marshal(metric)
}

func (s *service) ShutDown() {
	s.StoreMetrics()
}

func (s *service) RestoreMetrics() error {
	err := s.storage.RestoreMetrics()
	if err != nil {
		return err
	}
	return nil
}

func (s *service) StoreMetrics() error {
	err := s.storage.StoreMetrics()
	if err != nil {
		return err
	}
	return nil
}
