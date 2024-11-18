package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/memory"
	"github.com/ulixes-bloom/ya-metrics/internal/server/storage/pg"
)

func main() {
	// Инициализация конфигурации
	conf := config.Parse()
	// Инициализация контекста
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logLvl, err := zerolog.ParseLevel(conf.LogLvl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse log level")
	}
	zerolog.SetGlobalLevel(logLvl)

	// Инициализация хранилища метрик
	var storage service.Storage
	if conf.DatabaseDSN != "" {
		ps, err := pg.NewStorage(conf.DatabaseDSN)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		storage = ps
	} else {
		ms, err := memory.NewStorage(conf)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		storage = ms
	}

	err = api.New(conf, storage).Run(ctx)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
