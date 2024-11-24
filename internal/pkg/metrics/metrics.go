package metrics

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
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
