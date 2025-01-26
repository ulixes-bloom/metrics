package errors

import "errors"

var (
	ErrMetricNotExists          = errors.New("metric not exists")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricValueNotValid      = errors.New("metric value not valid")
)
