package memory

import (
	"bytes"
	"html/template"
	"strconv"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
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
	Metrics map[string]metrics.Metric
}

func NewMemStorage() *MemStorage {
	m := MemStorage{}
	m.Metrics = make(map[string]metrics.Metric,
		len(metrics.GaugeMetrics)+len(metrics.CounterMetrics))

	for _, g := range metrics.GaugeMetrics {
		zeroVal := float64(0)
		m.Metrics[g] = metrics.Metric{
			ID:    g,
			MType: metrics.Gauge,
			Value: &zeroVal,
		}
	}
	for _, c := range metrics.CounterMetrics {
		zeroVal := int64(0)
		m.Metrics[c] = metrics.Metric{
			ID:    c,
			MType: metrics.Counter,
			Delta: &zeroVal,
		}
	}
	return &m
}

func (m *MemStorage) Add(metric metrics.Metric) (metrics.Metric, error) {
	switch metric.MType {
	case metrics.Counter:
		cur := m.Metrics[metric.ID]
		newDelta := (*metric.Delta + *cur.Delta)
		metric.Delta = &newDelta
		m.Metrics[metric.ID] = metric
	case metrics.Gauge:
		m.Metrics[metric.ID] = metric
	default:
		return metric, errors.ErrMetricTypeNotImplemented
	}

	return metric, nil
}

func (m *MemStorage) Get(name string) (metrics.Metric, bool) {
	metric, ok := m.Metrics[name]
	return metric, ok
}

func (m *MemStorage) All() map[string]string {
	res := make(map[string]string)
	for k, v := range m.Metrics {
		switch v.MType {
		case metrics.Counter:
			res[k] = strconv.FormatInt(*v.Delta, 10)
		case metrics.Gauge:
			res[k] = strconv.FormatFloat(*v.Value, 'f', -1, 64)
		}
	}
	return res
}

func (m *MemStorage) HTMLTable() ([]byte, error) {
	var wr bytes.Buffer
	tmpl, err := template.New("tmpl").Parse(HTMLTemplate)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&wr, m.All())
	if err != nil {
		return nil, err
	}

	res := wr.Bytes()
	return res, nil
}
