package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

// GetMetric handles the HTTP request to retrieve a metric by its type and name.
// It responds with the metric's value if found, or an error if not.
func (a *httpApi) GetMetric(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	mtype := chi.URLParam(req, "mtype")
	mname := chi.URLParam(req, "mname")
	if mtype == "" || mname == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	mval, err := a.service.GetMetric(ctx, mtype, mname)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	res.Header().Add(headers.ContentType, "text/plain")
	res.Write([]byte(mval))
}

// UpdateMetric handles the HTTP request to update a metric's value by its type, name, and value.
// It responds with a success status or an error if the update fails.
func (a *httpApi) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	mtype := chi.URLParam(req, "mtype")
	mname := chi.URLParam(req, "mname")
	mval := chi.URLParam(req, "mval")
	if mtype == "" || mname == "" || mval == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.service.UpdateMetric(ctx, mtype, mname, mval)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// UpdateMetrics handles the HTTP request to update multiple metrics at once.
// It expects a JSON array of metrics, decodes it, and updates them in the service layer.
func (a *httpApi) UpdateMetrics(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var m []metrics.Metric
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&m); err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err := a.service.UpdateMetrics(ctx, m)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Add(headers.ContentType, "application/json")
	res.WriteHeader(http.StatusOK)
}

// GetMetricsHTMLTable handles the HTTP request to retrieve all metrics as HTML table.
// It responds with the generated HTML table or an error if the retrieval fails.
func (a *httpApi) GetMetricsHTMLTable(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	table, err := a.service.GetMetricsHTMLTable(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Add(headers.ContentType, "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(table)
}

// GetJSONMetric handles the HTTP request to retrieve a metric as a JSON object based on the provided metric name.
// It responds with the metric data in JSON format or an error if the metric is not found.
func (a *httpApi) GetJSONMetric(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var m metrics.Metric
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&m); err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	metric, err := a.service.GetJSONMetric(ctx, m)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	res.Header().Add(headers.ContentType, "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(metric)
}

// UpdateJSONMetric handles the HTTP request to update a metric based on the provided metric JSON object.
// It responds with the updated metric data in JSON format or an error if the update fails.
func (a *httpApi) UpdateJSONMetric(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var m metrics.Metric
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&m); err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := a.service.UpdateJSONMetric(ctx, m)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Add(headers.ContentType, "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(metric)
}

// PingDB handles the HTTP request to check the database connection status.
// It responds with a success status if the database is reachable, or an error if it is not.
func (a *httpApi) PingDB(res http.ResponseWriter, req *http.Request) {
	err := a.service.PingDB(a.conf.DatabaseDSN)
	if err != nil {
		log.Error().Msg(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
