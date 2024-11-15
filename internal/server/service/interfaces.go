package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type (
	// Интерфейс хранилища метрик
	Storage interface {
		Getter
		Setter

		Shutdown() error
	}

	// Интерфейс для получения метрик
	Getter interface {
		Get(name string) (val metrics.Metric, err error)
		GetAll() ([]metrics.Metric, error)
	}

	// Интерфейс для установки метрик
	Setter interface {
		Set(metric metrics.Metric) (metrics.Metric, error)
		SetAll(meticsSlice []metrics.Metric) error
	}
)
