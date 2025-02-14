package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
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

func (s *service) GetMetricsHTMLTable(ctx context.Context) ([]byte, error) {
	var wr bytes.Buffer
	tmpl, err := template.New("tmpl").Parse(metrics.HTMLTemplate)
	if err != nil {
		return nil, fmt.Errorf("service.getMetricsHTMLTable: %w", err)
	}
	allMetrics, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.getMetricsHTMLTable: %w", err)
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
		return nil, fmt.Errorf("service.getMetricsHTMLTable: %w", err)
	}

	res := wr.Bytes()
	return res, nil
}

func (s *service) GetMetric(ctx context.Context, mtype, mname string) ([]byte, error) {
	var mval string

	switch mtype {
	case metrics.Gauge:
		metric, err := s.storage.Get(ctx, mname)
		if err != nil {
			return nil, fmt.Errorf("service.getMetric: %w", err)
		}
		mval = strconv.FormatFloat(metric.GetValue(), 'f', -1, 64)
	case metrics.Counter:
		metric, err := s.storage.Get(ctx, mname)
		if err != nil {
			return nil, fmt.Errorf("service.getMetric: %w", err)
		}
		mval = strconv.FormatInt(metric.GetDelta(), 10)
	default:
		return nil, errors.ErrMetricTypeNotImplemented
	}

	return []byte(mval), nil
}

func (s *service) UpdateMetric(ctx context.Context, mtype, mname, mval string) error {
	switch mtype {
	case metrics.Gauge:
		if val, err := strconv.ParseFloat(mval, 64); err == nil {
			_, err := s.storage.Set(ctx, metrics.NewGaugeMetric(mname, val))
			if err != nil {
				return fmt.Errorf("service.updateMetric.parseGauge: %w", err)
			}
		} else {
			return errors.ErrMetricValueNotValid
		}
	case metrics.Counter:
		if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
			_, err := s.storage.Set(ctx, metrics.NewCounterMetric(mname, val))
			if err != nil {
				return fmt.Errorf("service.updateMetric.parseCounter: %w", err)
			}
		} else {
			return errors.ErrMetricValueNotValid
		}
	default:
		return errors.ErrMetricTypeNotImplemented
	}

	return nil
}

func (s *service) UpdateMetrics(ctx context.Context, metricsSlice []metrics.Metric) error {
	for _, m := range metricsSlice {
		_, err := s.UpdateJSONMetric(ctx, m)

		if err != nil {
			return fmt.Errorf("service.updateMetrics: %w", err)
		}
	}
	return nil
}

func (s *service) GetJSONMetric(ctx context.Context, metric metrics.Metric) ([]byte, error) {
	val, err := s.storage.Get(ctx, metric.ID)
	if err != nil {
		return nil, fmt.Errorf("service.getJSONMetric: %w", err)
	}

	return json.Marshal(val)
}

func (s *service) UpdateJSONMetric(ctx context.Context, metric metrics.Metric) ([]byte, error) {
	metric, err := s.storage.Set(ctx, metric)
	if err != nil {
		return nil, fmt.Errorf("service.updateJSONMetric: %w", err)
	}
	return json.Marshal(metric)
}

func (s *service) Shutdown(ctx context.Context) error {
	return s.storage.Shutdown(ctx)
}

func (s *service) PingDB(dsn string) error {
	err := retry.Do(func() error { return pg.PingDB(dsn) }, pg.NeedToRetry(), 4)
	if err != nil {
		return fmt.Errorf("service.pingDB: %w", err)
	}

	return nil
}
