package api

type Service interface {
	GetMetricsHTMLTable() ([]byte, error)
	GetMetric(mtype, mname string) ([]byte, error)
	UpdateMetric(mtype, mname, mval string) error
}
