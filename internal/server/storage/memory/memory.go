package memory

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"
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

type storage struct {
	metrics       map[string]metrics.Metric
	storeFilePath string
}

func NewStorage(storeFilePath string) *storage {
	s := storage{storeFilePath: storeFilePath}
	s.metrics = make(map[string]metrics.Metric,
		len(metrics.GaugeMetrics)+len(metrics.CounterMetrics))
	for _, g := range metrics.GaugeMetrics {
		zeroVal := float64(0)
		s.metrics[g] = metrics.Metric{
			ID:    g,
			MType: metrics.Gauge,
			Value: &zeroVal,
		}
	}
	for _, c := range metrics.CounterMetrics {
		zeroVal := int64(0)
		s.metrics[c] = metrics.Metric{
			ID:    c,
			MType: metrics.Counter,
			Delta: &zeroVal,
		}
	}
	return &s
}

func (s *storage) Add(metric metrics.Metric) (metrics.Metric, error) {
	switch metric.MType {
	case metrics.Counter:
		cur, ok := s.metrics[metric.ID]
		if ok {
			newDelta := (*metric.Delta + *cur.Delta)
			metric.Delta = &newDelta
			s.metrics[metric.ID] = metric
		} else {
			s.metrics[metric.ID] = metric
		}
	case metrics.Gauge:
		s.metrics[metric.ID] = metric
	default:
		return metric, errors.ErrMetricTypeNotImplemented
	}

	return metric, nil
}

func (s *storage) Get(name string) (metrics.Metric, bool) {
	metric, ok := s.metrics[name]
	return metric, ok
}

func (s *storage) All() map[string]string {
	res := make(map[string]string)
	for k, v := range s.metrics {
		switch v.MType {
		case metrics.Counter:
			res[k] = strconv.FormatInt(*v.Delta, 10)
		case metrics.Gauge:
			res[k] = strconv.FormatFloat(*v.Value, 'f', -1, 64)
		}
	}
	return res
}

func (s *storage) HTMLTable() ([]byte, error) {
	var wr bytes.Buffer
	tmpl, err := template.New("tmpl").Parse(HTMLTemplate)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&wr, s.All())
	if err != nil {
		return nil, err
	}

	res := wr.Bytes()
	return res, nil
}

func (s *storage) RestoreMetrics() error {
	file, err := os.OpenFile(s.storeFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	var restoredMetrics map[string]metrics.Metric
	err = json.NewDecoder(file).Decode(&restoredMetrics)
	if err != nil {
		return err
	}

	s.metrics = restoredMetrics
	return nil
}

func (s *storage) StoreMetrics() error {
	file, err := os.OpenFile(s.storeFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	dataMetricsJSON, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	_, err = file.Write(dataMetricsJSON)
	if err != nil {
		return err
	}

	return nil
}
