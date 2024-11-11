package api

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/logger"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/memory"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/pg"
)

type api struct {
	service Service
	conf    *config.Config
	log     zerolog.Logger
	router  *chi.Mux
}

func New(conf *config.Config) *api {
	// Инициализация логгера
	logger, err := logger.Initialize(conf.LogLvl, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация хранилища метрик
	var storage service.Storage
	if conf.DatabaseDSN != "" {
		ps, err := pg.NewStorage(conf.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		storage = ps
	} else {
		storage = memory.NewStorage(logger, conf)
	}

	srv := service.New(storage, logger, conf)
	newAPI := api{
		service: srv,
		log:     logger,
		conf:    conf,
	}
	newAPI.router = newAPI.newRouter()
	return &newAPI
}

func (a *api) Run(ctx context.Context) {
	go func() {
		err := http.ListenAndServe(a.conf.RunAddr, a.router)
		if err != nil {
			a.log.Fatal().Msg(err.Error())
		}
	}()

	<-ctx.Done()
	a.service.Shutdown()
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
	r.Post("/updates/", a.UpdateMetrics)
	return r
}
