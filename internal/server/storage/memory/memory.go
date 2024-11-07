package memory

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
)

type memstorage struct {
	metrics map[string]metrics.Metric
	log     zerolog.Logger
	conf    config.Config
}

func NewStorage(logger zerolog.Logger, conf config.Config) *memstorage {
	ms := memstorage{
		log:  logger,
		conf: conf,
	}
	ms.metrics = make(map[string]metrics.Metric,
		len(metrics.GaugeMetrics)+len(metrics.CounterMetrics))
	for _, g := range metrics.GaugeMetrics {
		zeroVal := float64(0)
		ms.metrics[g] = metrics.Metric{
			ID:    g,
			MType: metrics.Gauge,
			Value: &zeroVal,
		}
	}
	for _, c := range metrics.CounterMetrics {
		zeroVal := int64(0)
		ms.metrics[c] = metrics.Metric{
			ID:    c,
			MType: metrics.Counter,
			Delta: &zeroVal,
		}
	}
	return &ms
}

func (ms *memstorage) Set(metric metrics.Metric) (metrics.Metric, error) {
	switch metric.MType {
	case metrics.Counter:
		cur, ok := ms.metrics[metric.ID]
		if ok {
			newDelta := (*metric.Delta + *cur.Delta)
			metric.Delta = &newDelta
			ms.metrics[metric.ID] = metric
		} else {
			ms.metrics[metric.ID] = metric
		}
	case metrics.Gauge:
		ms.metrics[metric.ID] = metric
	default:
		return metric, errors.ErrMetricTypeNotImplemented
	}

	return metric, nil
}

func (ms *memstorage) Get(name string) (metrics.Metric, bool) {
	metric, ok := ms.metrics[name]
	return metric, ok
}

func (ms *memstorage) GetAll() ([]metrics.Metric, error) {
	allMetrics := make([]metrics.Metric, 0)
	for _, m := range ms.metrics {
		allMetrics = append(allMetrics, m)
	}
	return allMetrics, nil
}

func (ms *memstorage) Setup() error {
	err := ms.restoreMetricsFromFile()
	ms.async()
	return err
}

func (ms *memstorage) Shutdown() error {
	return ms.saveMetricsToFile()
}

func (ms *memstorage) async() {
	if ms.conf.StoreInterval == 0 {
		return
	}
	storeTicker := time.NewTicker(ms.conf.StoreInterval)

	go func() {
		for {
			<-storeTicker.C

			if err := ms.saveMetricsToFile(); err != nil {
				ms.log.Err(err)
			}
		}
	}()
}

func (ms *memstorage) restoreMetricsFromFile() error {
	file, err := os.OpenFile(ms.conf.FileStoragePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	var restoredMetrics map[string]metrics.Metric
	err = json.NewDecoder(file).Decode(&restoredMetrics)
	if err != nil {
		return err
	}
	ms.metrics = restoredMetrics
	return nil
}

func (ms *memstorage) saveMetricsToFile() error {
	file, err := os.OpenFile(ms.conf.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	if err = encoder.Encode(ms.metrics); err != nil {
		return err
	}
	return writer.Flush()
}
