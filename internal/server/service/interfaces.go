package service

import "github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"

type (
	// Интерфейс хранилища метрик
	Storage interface {
		Getter
		Setter

		Setup() error
		Shutdown() error
	}

	// Интерфейс для получения метрик
	Getter interface {
		Get(name string) (val metrics.Metric, ok bool)
		GetAll() ([]metrics.Metric, error)
	}

	// Интерфейс для установки метрик
	Setter interface {
		Set(metric metrics.Metric) (metrics.Metric, error)
		SetAll(meticsSlice []metrics.Metric) error
	}
)
