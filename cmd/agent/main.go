package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/client"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/memory"
)

func main() {
	// Инициализация конфигурации
	conf, err := config.Parse()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	// Инициализация контекста
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()
	// Инициализация логгера
	logLvl, err := zerolog.ParseLevel(conf.LogLvl)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse log level")
	}
	zerolog.SetGlobalLevel(logLvl)
	// Инициализация in memory хранилища метрик
	ms := memory.NewStorage()

	cl := client.New(conf, ms)
	cl.Run(ctx)
}
