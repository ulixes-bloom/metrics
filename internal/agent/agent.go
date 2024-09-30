package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/storage"
)

type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	MemStorage     *storage.MemStorage
	ServerAddr     string
}

func NewAgent(pollInterval, reportInterval time.Duration, serverAddr string) *Agent {
	return &Agent{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		MemStorage:     storage.NewMemStorage(),
		ServerAddr:     serverAddr,
	}
}

func (a *Agent) UpdateGuageMetrics() {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	reflectMemStats := reflect.ValueOf(ms)
	for _, metricName := range metrics.GaugeMetrics {
		reflectMemStatsField := reflectMemStats.FieldByName(metricName)

		if reflectMemStatsField.Kind() == reflect.Invalid {
			continue
		}

		var metricVal float64

		switch reflectMemStatsVal := reflectMemStatsField.Interface().(type) {
		case float64:
			metricVal = reflectMemStatsVal
		case uint64:
			metricVal = float64(reflectMemStatsVal)
		case uint32:
			metricVal = float64(reflectMemStatsVal)
		default:
			fmt.Println("unexpected metric type")
		}

		a.MemStorage.AddGauge(metricName, metricVal)
	}

	a.MemStorage.AddGauge("RandomValue", rand.Float64())
}

func (a *Agent) UpdateCounterMetrics() {
	a.MemStorage.AddCounter("PollCount", 1)
}

func (a *Agent) SendMetrics() {
	for k, v := range a.MemStorage.Gauges {
		a.SendMetric("gauge", k, fmt.Sprintf("%v", v))
	}

	for k, v := range a.MemStorage.Counters {
		a.SendMetric("counter", k, fmt.Sprintf("%d", v))
	}
}

func (a *Agent) SendMetric(mtype, mname, mval string) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", a.ServerAddr, mtype, mname, mval)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
	}
	defer resp.Body.Close()
}
