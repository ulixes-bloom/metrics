package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type (
	Storage interface {
		Getter
		Setter

		Shutdown() error
	}

	Getter interface {
		Get(name string) (val metrics.Metric, err error)
		GetAll() ([]metrics.Metric, error)
	}

	Setter interface {
		Set(metric metrics.Metric) (metrics.Metric, error)
		SetAll(meticsSlice []metrics.Metric) error
	}
)
