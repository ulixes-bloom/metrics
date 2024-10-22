package api

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type Service interface {
	GetMetricsHTMLTable() ([]byte, error)
	GetJSONMetric(mtype, mname string) ([]byte, error)
	UpdateJSONMetric(metric metrics.Metric) ([]byte, error)
}
