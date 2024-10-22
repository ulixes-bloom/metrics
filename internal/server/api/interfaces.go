package api

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Service interface {
	GetMetric(mtype, mname string) ([]byte, error)
	UpdateMetric(mtype, mname, mval string) error
	GetMetricsHTMLTable() ([]byte, error)
	GetJSONMetric(metric metrics.Metric) ([]byte, error)
	UpdateJSONMetric(metric metrics.Metric) ([]byte, error)
}
