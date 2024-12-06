package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api/middleware"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
)

type api struct {
	service Service
	conf    *config.Config
	router  *chi.Mux
}

func New(conf *config.Config, storage service.Storage) *api {
	srv := service.New(storage, conf)
	newAPI := api{
		service: srv,
		conf:    conf,
	}
	newAPI.router = newAPI.newRouter()
	return &newAPI
}

func (a *api) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- http.ListenAndServe(a.conf.RunAddr, a.router)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return a.service.Shutdown()
	}
}

func (a *api) newRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.WithLogging)
	if a.conf.HashKey != "" {
		r.Use(middleware.WithHashing(a.conf.HashKey))
	}
	r.Use(middleware.WithCompressing)

	r.Get("/", a.GetMetricsHTMLTable)
	r.Get("/ping", a.PingDB)
	r.Get("/value/{mtype}/{mname}", a.GetMetric)
	r.Post("/update/{mtype}/{mname}/{mval}", a.UpdateMetric)
	r.Post("/value/", a.GetJSONMetric)
	r.Post("/update/", a.UpdateJSONMetric)
	r.Post("/updates/", a.UpdateMetrics)
	return r
}
