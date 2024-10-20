package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Sevice Service
}

func NewHandler(service Service) *Handler {
	return &Handler{Sevice: service}
}

func (h *Handler) GetMetricsHTMLTable(res http.ResponseWriter, req *http.Request) {
	table, err := h.Sevice.GetMetricsHTMLTable()
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(table)
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
		log.Error().Msg(err.Error())
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
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.Sevice.UpdateMetric(mtype, mname, mval)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
