package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/storage"
)

var MemStorage = storage.NewMemStorage()

func GetMetricsList(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(MemStorage.HTMLTable())
}

func GetMetric(res http.ResponseWriter, req *http.Request) {
	mtype := chi.URLParam(req, "mtype")
	mname := chi.URLParam(req, "mname")
	if mtype == "" || mname == "" {
		res.WriteHeader(http.StatusBadRequest)
	}
	var mval string

	switch mtype {
	case "gauge":
		val, ok := MemStorage.GetGauge(mname)
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		mval = strconv.FormatFloat(val, 'f', -1, 64)
	case "counter":
		val, ok := MemStorage.GetCounter(mname)
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		mval = strconv.FormatInt(val, 10)
	default:
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Add("Content-Type", "text/plain")
	res.Write([]byte(mval))
}

func UpdateMetric(res http.ResponseWriter, req *http.Request) {
	mtype := chi.URLParam(req, "mtype")
	mname := chi.URLParam(req, "mname")
	mval := chi.URLParam(req, "mval")

	if mtype == "" || mname == "" || mval == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch mtype {
	case "gauge":
		if val, err := strconv.ParseFloat(mval, 64); err == nil {
			MemStorage.AddGauge(mname, val)
		} else {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	case "counter":
		if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
			MemStorage.AddCounter(mname, val)
		} else {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
