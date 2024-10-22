package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type Handler struct {
	Sevice Service
	Logger zerolog.Logger
}

func NewHandler(service Service, logger zerolog.Logger) *Handler {
	return &Handler{
		Sevice: service,
		Logger: logger,
	}
}

func (h *Handler) GetMetricsHTMLTable(res http.ResponseWriter, req *http.Request) {
	table, err := h.Sevice.GetMetricsHTMLTable()
	if err != nil {
		h.Logger.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(table)
}

func (h *Handler) GetJSONMetric(res http.ResponseWriter, req *http.Request) {
	var m metrics.Metric
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&m); err != nil {
		h.Logger.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	metric, err := h.Sevice.GetJSONMetric(m.MType, m.ID)
	if err != nil {
		h.Logger.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	res.Header().Add("Content-Type", "application/json")
	res.Write(metric)
}

func (h *Handler) UpdateJSONMetric(res http.ResponseWriter, req *http.Request) {
	var m metrics.Metric
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&m); err != nil {
		h.Logger.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Sevice.UpdateJSONMetric(m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(res)
	if err := enc.Encode(metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}
