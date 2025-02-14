package service

import (
	"context"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type (
	Storage interface {
		Getter
		Setter

		Shutdown(ctx context.Context) error
	}

	Getter interface {
		Get(ctx context.Context, name string) (val metrics.Metric, err error)
		GetAll(ctx context.Context) ([]metrics.Metric, error)
	}

	Setter interface {
		Set(ctx context.Context, metric metrics.Metric) (metrics.Metric, error)
		SetAll(ctx context.Context, meticsSlice []metrics.Metric) error
	}
)
