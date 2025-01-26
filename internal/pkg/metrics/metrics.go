// Package metrics provides functionality for defining and operating with metrics,
// including both "gauge" and "counter" types, and supports serialization to JSON.
package metrics

// Metric represents a single metric with an ID, type, and a value or delta depending on the metric type.
type Metric struct {
	ID    string   `json:"id"`              // Metric name (ID)
	MType string   `json:"type"`            // Metric type: "gauge" or "counter"
	Delta *int64   `json:"delta,omitempty"` // Delta value for counter type metrics (optional)
	Value *float64 `json:"value,omitempty"` // Value for gauge type metrics (optional)
}

func NewGaugeMetric(id string, val float64) Metric {
	return Metric{
		ID:    id,
		MType: Gauge,
		Value: &val,
	}
}

func NewCounterMetric(id string, delta int64) Metric {
	return Metric{
		ID:    id,
		MType: Counter,
		Delta: &delta,
	}
}

func (m *Metric) GetDelta() int64 {
	if m.Delta == nil {
		return 0
	}
	return *m.Delta
}

func (m *Metric) GetValue() float64 {
	if m.Value == nil {
		return 0
	}
	return *m.Value
}

const Counter = "counter"
const Gauge = "gauge"

var (
	CounterMetrics = []string{
		"PollCount",
	}

	GaugeMetrics = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
		"TotalMemory",
		"FreeMemory",
		"CPUutilization1",
	}

	MetricsCount = len(GaugeMetrics) + len(CounterMetrics)
)

// HTMLTemplate is the HTML template used to render the metrics in a web page.
const HTMLTemplate = `<html>
	<head>
		<title></title>
	</head>
	<body>
		<table>
			<tr>
				<th>Метрика</th>
				<th>Значение</th>
			</tr>
			{{range $key, $value := .}}
			<tr>
				<td>{{$key}}</td>
				<td>{{$value}}</td>
			</tr>
			{{end}}
		</table>
	</body>
</html>`
