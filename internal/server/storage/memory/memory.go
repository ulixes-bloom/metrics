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
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metricerrors"
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
	ms.metrics = map[string]metrics.Metric{}
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

	err := ms.setup()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("memory.newStorage.setup: %w", err)
	}

	return &ms, nil
}

func (ms *memstorage) Set(metric metrics.Metric) (metrics.Metric, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	switch metric.MType {
	case metrics.Counter:
		cur, ok := ms.metrics[metric.ID]
		if ok {
			newDelta := (metric.GetDelta() + cur.GetDelta())
			metric.Delta = &newDelta
			ms.metrics[metric.ID] = metric
		} else {
			ms.metrics[metric.ID] = metric
		}
	case metrics.Gauge:
		ms.metrics[metric.ID] = metric
	default:
		return metric, metricerrors.ErrMetricTypeNotImplemented
	}

	if err := ms.sync(); err != nil {
		return metric, fmt.Errorf("memory.set.sync: %w", err)
	}
	return metric, nil
}

func (ms *memstorage) SetAll(meticsSlice []metrics.Metric) error {
	for _, m := range meticsSlice {
		if _, err := ms.Set(m); err != nil {
			return fmt.Errorf("memory.setAll.set: %w", err)
		}
	}

	if err := ms.sync(); err != nil {
		return fmt.Errorf("memory.setAll.sync: %w", err)
	}
	return nil
}

func (ms *memstorage) Get(name string) (metrics.Metric, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	metric, ok := ms.metrics[name]
	if !ok {
		return metric, metricerrors.ErrMetricNotExists
	}
	return metric, nil
}

func (ms *memstorage) GetAll() ([]metrics.Metric, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	allMetrics := []metrics.Metric{}
	for _, m := range ms.metrics {
		allMetrics = append(allMetrics, m)
	}
	return allMetrics, nil
}

func (ms *memstorage) setup() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	err := ms.restoreMetricsFromFile()
	if err != nil {
		return fmt.Errorf("memory.setup: %w", err)
	}

	ms.async()
	return nil
}

func (ms *memstorage) Shutdown() error {
	return ms.saveMetricsToFile()
}

// Асинхронная запись метрик в файл
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

// Синхронная запись метрик в файл
func (ms *memstorage) sync() error {
	if ms.conf.StoreInterval == 0 {
		err := ms.saveMetricsToFile()
		if err != nil {
			return fmt.Errorf("memory.sync: %w", err)
		}
	}
	return nil
}

// Считывание значений метрик из файла в память
func (ms *memstorage) restoreMetricsFromFile() error {
	_, error := os.Stat(ms.conf.FileStoragePath)
	if errors.Is(error, os.ErrNotExist) {
		return nil
	}

	file, err := os.OpenFile(ms.conf.FileStoragePath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("memory.restoreMetricsFromFile.openFile: '%s', %w", ms.conf.FileStoragePath, err)
	}
	var restoredMetrics map[string]metrics.Metric
	err = json.NewDecoder(file).Decode(&restoredMetrics)
	if err != nil {
		return fmt.Errorf("memory.restoreMetricsFromFile.decode: %w", err)
	}
	ms.metrics = restoredMetrics
	return nil
}

// Сохранение значений метрик из памяти в файл
func (ms *memstorage) saveMetricsToFile() error {
	file, err := os.OpenFile(ms.conf.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("memory.saveMetricsToFile.openFile: '%s', %w", ms.conf.FileStoragePath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	msMetrics, _ := ms.GetAll()
	if err = encoder.Encode(msMetrics); err != nil {
		return fmt.Errorf("memory.saveMetricsToFile.encode: %w", err)
	}
	return writer.Flush()
}
