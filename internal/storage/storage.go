package storage

import (
	"bytes"
	"html/template"
	"log"
	"strconv"

	"github.com/ulixes-bloom/ya-metrics/internal/metrics"
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

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMemStorage() *MemStorage {
	m := MemStorage{}
	m.Gauges = make(map[string]float64, len(metrics.GaugeMetrics))
	m.Counters = make(map[string]int64, len(metrics.CounterMetrics))
	return &m
}

func (m *MemStorage) AddGauge(name string, value float64) {
	m.Gauges[name] = value
}

func (m *MemStorage) AddCounter(name string, value int64) {
	m.Counters[name] += value
}

func (m *MemStorage) GetGauge(name string) (val float64, ok bool) {
	val, ok = m.Gauges[name]
	return
}

func (m *MemStorage) GetCounter(name string) (val int64, ok bool) {
	val, ok = m.Counters[name]
	return
}

func (m *MemStorage) All() map[string]string {
	res := make(map[string]string)

	for k, v := range m.Gauges {
		res[k] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	for k, v := range m.Counters {
		res[k] = strconv.FormatInt(v, 10)
	}

	return res
}

func (m *MemStorage) HTMLTable() (res []byte) {
	var wr bytes.Buffer
	tmpl, err := template.New("tmpl").Parse(HTMLTemplate)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&wr, m.All())
	if err != nil {
		log.Fatal(err)
	}
	res = wr.Bytes()
	return
}
