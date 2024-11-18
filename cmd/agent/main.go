package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/client"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
)

func main() {
	// Инициализация конфигурации
	conf, err := config.Parse()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	// Инициализация контекста
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logLvl, err := zerolog.ParseLevel(conf.LogLvl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse log level")
	}
	zerolog.SetGlobalLevel(logLvl)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	client.New(conf).Run(ctx)
}
