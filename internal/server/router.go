package server

import "github.com/go-chi/chi/v5"

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", GetMetricsList)
	r.Get("/value/{mtype}/{mname}", GetMetric)
	r.Post("/update/{mtype}/{mname}/{mval}", UpdateMetric)

	return r
}
