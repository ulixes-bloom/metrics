package api

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/logger"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/memory"
)

func Run(addr, loglvl string) error {
	r := Router(loglvl)
	err := http.ListenAndServe(addr, r)
	return err
}

func Router(loglvl string) chi.Router {
	st := memory.NewMemStorage()
	srv := service.NewService(st)

	logger, err := logger.Initialize(loglvl, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	h := NewHandler(srv, logger)

	r := chi.NewRouter()
	r.Use(h.WithLogging)
	r.Use(h.WithCompressing)
	r.Get("/", h.GetMetricsHTMLTable)
	r.Get("/value/{mtype}/{mname}", h.GetMetric)
	r.Post("/update/{mtype}/{mname}/{mval}", h.UpdateMetric)
	r.Post("/value/", h.GetJSONMetric)
	r.Post("/update/", h.UpdateJSONMetric)
	return r
}
