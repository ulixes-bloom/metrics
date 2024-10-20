package service

type Storage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (val float64, ok bool)
	GetCounter(name string) (val int64, ok bool)
	All() map[string]string
	HTMLTable() ([]byte, error)
}
