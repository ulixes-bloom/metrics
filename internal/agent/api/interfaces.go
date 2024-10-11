package api

type Service interface {
	UpdateMetrics()
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}
