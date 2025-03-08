package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api/middleware"
)

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

	r.Group(func(r chi.Router) {
		if a.conf.TrustedSubnet != "" {
			r.Use(middleware.WithIPResolving(a.conf.TrustedSubnet))
		}
		if a.conf.CryptoKey != "" {
			r.Use(middleware.WithRSA(a.conf.CryptoKey))
		}
		r.Post("/value/", a.GetJSONMetric)
		r.Post("/update/", a.UpdateJSONMetric)
		r.Post("/updates/", a.UpdateMetrics)
	})
	return r
}
