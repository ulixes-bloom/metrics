package metrics

type MemStorage struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

func NewMemStorage() *MemStorage {
	m := MemStorage{}
	m.Counters = make(map[string]Counter)
	m.Gauges = make(map[string]Gauge)
	return &m
}

func (m *MemStorage) AddGauge(name string, value Gauge) {
	m.Gauges[name] = value
}

func (m *MemStorage) AddCounter(name string, value Counter) {
	m.Counters[name] += value
}

type Gauge float64

type Counter int64
