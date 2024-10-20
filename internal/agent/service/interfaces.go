package service

type Storage interface {
	Add(name string, value interface{}) error
	GetGauge(name string) (val float64, ok bool)
	GetCounter(name string) (val int64, ok bool)
	GetAllCounters() map[string]int64
	GetAllGauges() map[string]float64
}
