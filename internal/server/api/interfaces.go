package api

import (
	"context"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type Service interface {
	GetMetric(ctx context.Context, mtype, mname string) ([]byte, error)
	UpdateMetric(ctx context.Context, mtype, mname, mval string) error
	UpdateMetrics(ctx context.Context, m []metrics.Metric) error
	GetMetricsHTMLTable(ctx context.Context) ([]byte, error)
	GetJSONMetric(ctx context.Context, metric metrics.Metric) ([]byte, error)
	UpdateJSONMetric(ctx context.Context, metric metrics.Metric) ([]byte, error)
	PingDB(dsn string) error
	Shutdown(ctx context.Context) error
}
