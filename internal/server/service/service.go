package service

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strconv"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metricerrors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/retry"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/pg"
)

type service struct {
	storage Storage
	conf    *config.Config
}

func New(storage Storage, conf *config.Config) *service {
	srv := &service{
		storage: storage,
		conf:    conf,
	}

	return srv
}

func (s *service) GetMetricsHTMLTable() ([]byte, error) {
	var wr bytes.Buffer
	tmpl, err := template.New("tmpl").Parse(metrics.HTMLTemplate)
	if err != nil {
		return nil, err
	}
	allMetrics, err := s.storage.GetAll()
	if err != nil {
		return nil, err
	}
	metricsMap := map[string]string{}

	for _, m := range allMetrics {
		switch m.MType {
		case metrics.Counter:
			metricsMap[m.ID] = strconv.FormatInt(m.GetDelta(), 10)
		case metrics.Gauge:
			metricsMap[m.ID] = strconv.FormatFloat(m.GetValue(), 'f', -1, 64)
		}
	}

	err = tmpl.Execute(&wr, metricsMap)
	if err != nil {
		return nil, err
	}

	res := wr.Bytes()
	return res, nil
}

func (s *service) GetMetric(mtype, mname string) ([]byte, error) {
	var mval string

	switch mtype {
	case metrics.Gauge:
		metric, err := s.storage.Get(mname)
		if err != nil {
			return nil, err
		}
		mval = strconv.FormatFloat(metric.GetValue(), 'f', -1, 64)
	case metrics.Counter:
		metric, err := s.storage.Get(mname)
		if err != nil {
			return nil, err
		}
		mval = strconv.FormatInt(metric.GetDelta(), 10)
	default:
		return nil, metricerrors.ErrMetricTypeNotImplemented
	}

	return []byte(mval), nil
}

func (s *service) UpdateMetric(mtype, mname, mval string) error {
	switch mtype {
	case metrics.Gauge:
		if val, err := strconv.ParseFloat(mval, 64); err == nil {
			_, err := s.storage.Set(metrics.NewGaugeMetric(mname, val))
			if err != nil {
				return err
			}
		} else {
			return metricerrors.ErrMetricValueNotValid
		}
	case metrics.Counter:
		if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
			_, err := s.storage.Set(metrics.NewCounterMetric(mname, val))
			if err != nil {
				return err
			}
		} else {
			return metricerrors.ErrMetricValueNotValid
		}
	default:
		return metricerrors.ErrMetricTypeNotImplemented
	}

	return nil
}

func (s *service) UpdateMetrics(metricsSlice []metrics.Metric) error {
	for _, m := range metricsSlice {
		_, err := s.UpdateJSONMetric(m)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetJSONMetric(metric metrics.Metric) ([]byte, error) {
	val, err := s.storage.Get(metric.ID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(val)
}

func (s *service) UpdateJSONMetric(metric metrics.Metric) ([]byte, error) {
	metric, err := s.storage.Set(metric)
	if err != nil {
		return nil, err
	}
	return json.Marshal(metric)
}

func (s *service) Shutdown() error {
	return s.storage.Shutdown()
}

func (s *service) PingDB(dsn string) error {
	err := retry.Do(func() error { return pg.PingDB(dsn) }, pg.NeedToRetry(), 4)
	if err != nil {
		return err
	}

	return nil
}
