package metricerrors

import "errors"

var (
	ErrMetricNotExists          = errors.New("metric not exists")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricValueNotValid      = errors.New("metric value not valid")
	ErrFailedMetricCompression  = errors.New("failed to compress metric")
	ErrFailedMetricMarshall     = errors.New("failed to marshall metric")
)
