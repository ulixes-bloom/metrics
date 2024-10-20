package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/memory"
)

func Run(addr string) error {
	r := Router()
	err := http.ListenAndServe(addr, r)
	return err
}

func Router() chi.Router {
	st := memory.NewMemStorage()
	srv := service.NewService(st)
	h := NewHandler(srv)

	r := chi.NewRouter()
	r.Get("/", h.GetMetricsHTMLTable)
	r.Get("/value/{mtype}/{mname}", h.GetMetric)
	r.Post("/update/{mtype}/{mname}/{mval}", h.UpdateMetric)

	return r
}
