package memory

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	appErrors "github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
)

type memstorage struct {
	metrics map[string]metrics.Metric
	conf    *config.Config
	mutex   sync.RWMutex
}

func NewStorage(conf *config.Config) (*memstorage, error) {
	ms := memstorage{
		conf: conf,
	}
	// pre-allocate the metrics map with the expected size
	ms.metrics = make(map[string]metrics.Metric, metrics.MetricsCount)

	// initialize Gauge metrics with a default value of 0
	for _, g := range metrics.GaugeMetrics {
		zeroVal := float64(0)
		ms.metrics[g] = metrics.Metric{
			ID:    g,
			MType: metrics.Gauge,
			Value: &zeroVal,
		}
	}

	// initialize Counter metrics with a default value of 0
	for _, c := range metrics.CounterMetrics {
		zeroVal := int64(0)
		ms.metrics[c] = metrics.Metric{
			ID:    c,
			MType: metrics.Counter,
			Delta: &zeroVal,
		}
	}

	err := ms.setup()
	if err != nil {
		return nil, fmt.Errorf("memory.newStorage.setup: %w", err)
	}

	return &ms, nil
}

func (ms *memstorage) Set(metric metrics.Metric) (metrics.Metric, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	switch metric.MType {
	case metrics.Counter:
		cur, exists := ms.metrics[metric.ID]
		if exists {
			newDelta := metric.GetDelta() + cur.GetDelta()
			metric.Delta = &newDelta
		}
		ms.metrics[metric.ID] = metric
	case metrics.Gauge:
		ms.metrics[metric.ID] = metric
	default:
		return metric, appErrors.ErrMetricTypeNotImplemented
	}

	if err := ms.sync(); err != nil {
		return metric, fmt.Errorf("memory.set: %w", err)
	}
	return metric, nil
}

func (ms *memstorage) SetAll(metricsSlice []metrics.Metric) error {
	for _, m := range metricsSlice {
		if _, err := ms.Set(m); err != nil {
			return fmt.Errorf("memory.setAll.set: %w", err)
		}
	}

	if err := ms.sync(); err != nil {
		return fmt.Errorf("memory.setAll: %w", err)
	}
	return nil
}

func (ms *memstorage) Get(name string) (metrics.Metric, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	metric, exists := ms.metrics[name]
	if !exists {
		return metric, appErrors.ErrMetricNotExists
	}
	return metric, nil
}

func (ms *memstorage) GetAll() ([]metrics.Metric, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	allMetrics := make([]metrics.Metric, 0, len(ms.metrics))
	for _, m := range ms.metrics {
		allMetrics = append(allMetrics, m)
	}
	return allMetrics, nil
}

func (ms *memstorage) setup() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Restore metrics values from file
	if err := ms.restoreMetricsFromFile(); err != nil {
		return fmt.Errorf("memory.setup: %w", err)
	}

	ms.async()
	return nil
}

func (ms *memstorage) Shutdown() error {
	return ms.saveMetricsToFile()
}

// start a background process to save metrics to a file with period conf.StoreInterval
func (ms *memstorage) async() {
	if ms.conf.StoreInterval == 0 {
		return
	}
	storeTicker := time.NewTicker(ms.conf.StoreInterval)

	go func() {
		for {
			<-storeTicker.C
			if err := ms.saveMetricsToFile(); err != nil {
				log.Err(err)
			}
		}
	}()
}

// persist metrics values to a file if config.StoreInterval is set to 0.
func (ms *memstorage) sync() error {
	if ms.conf.StoreInterval == 0 {
		if err := ms.saveMetricsToFile(); err != nil {
			return fmt.Errorf("memory.sync: %w", err)
		}
	}
	return nil
}

// load metrics from the file into memory.
func (ms *memstorage) restoreMetricsFromFile() error {
	file, err := os.OpenFile(ms.conf.FileStoragePath, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("memory.restoreMetricsFromFile.openFile: '%s', %w", ms.conf.FileStoragePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Debug().Msgf("memory.restoreMetricsFromFile: failed to close file '%s'", ms.conf.FileStoragePath)
		}
	}()

	restoredMetrics := make(map[string]metrics.Metric)
	if err = json.NewDecoder(file).Decode(&restoredMetrics); err != nil {
		return fmt.Errorf("memory.restoreMetricsFromFile.decode: %w", err)
	}
	ms.metrics = restoredMetrics
	return nil
}

// write metrics from memory to the file specified in config.FileStoragePath
func (ms *memstorage) saveMetricsToFile() error {
	file, err := os.OpenFile(ms.conf.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("memory.saveMetricsToFile.openFile: '%s', %w", ms.conf.FileStoragePath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	msMetrics, _ := ms.GetAll()
	if err := encoder.Encode(msMetrics); err != nil {
		return fmt.Errorf("memory.saveMetricsToFile.encode: %w", err)
	}
	return writer.Flush()
}
