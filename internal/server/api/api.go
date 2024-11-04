package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/logger"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/memory"
)

type api struct {
	service Service
	config  config.Config
	log     zerolog.Logger
	router  *chi.Mux
}

func New(conf config.Config) *api {
	logger, err := logger.Initialize(conf.LogLvl, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	st := memory.NewStorage(conf.FileStoragePath)
	srv := service.New(st, logger, conf)

	newAPI := api{
		service: srv,
		log:     logger,
		config:  conf,
	}
	newAPI.router = newAPI.newRouter()

	return &newAPI
}

func (a *api) Run(ctx context.Context) {
	go func() {
		err := http.ListenAndServe(a.config.RunAddr, a.router)
		if err != nil {
			a.log.Fatal().Msg(err.Error())
		}
	}()

	storeTicker := time.NewTicker(a.config.StoreInterval)
	defer storeTicker.Stop()

	for {
		select {
		case <-storeTicker.C:
			a.service.StoreMetrics()
		case <-ctx.Done():
			a.service.ShutDown()
			return
		}
	}
}

func (a *api) newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(a.MiddlewareLogging)
	r.Use(a.MiddlewareCompressing)
	r.Get("/", a.GetMetricsHTMLTable)
	r.Get("/ping", a.PingDB)
	r.Get("/value/{mtype}/{mname}", a.GetMetric)
	r.Post("/update/{mtype}/{mname}/{mval}", a.UpdateMetric)
	r.Post("/value/", a.GetJSONMetric)
	r.Post("/update/", a.UpdateJSONMetric)
	return r
}
