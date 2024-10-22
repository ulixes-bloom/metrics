package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) GetMetric(res http.ResponseWriter, req *http.Request) {
	mtype := chi.URLParam(req, "mtype")
	mname := chi.URLParam(req, "mname")
	if mtype == "" || mname == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	mval, err := h.Sevice.GetMetric(mtype, mname)
	if err != nil {
		h.Logger.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	res.Header().Add("Content-Type", "text/plain")
	res.Write([]byte(mval))
}

func (h *Handler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	mtype := chi.URLParam(req, "mtype")
	mname := chi.URLParam(req, "mname")
	mval := chi.URLParam(req, "mval")
	if mtype == "" || mname == "" || mval == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Sevice.UpdateMetric(mtype, mname, mval)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
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

	metric, err := h.Sevice.GetJSONMetric(m)
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
		h.Logger.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.Write(metric)
}
