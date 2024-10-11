package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/api"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/storage"
)

func main() {
	conf, err := parseConfig()

	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	m := storage.NewMemStorage()
	s := service.NewService(m)
	a := api.NewClient(s,
		time.Duration(conf.PollInterval)*time.Second,
		time.Duration(conf.ReportInterval)*time.Second,
		"http://"+conf.ServerAddr)

	go func() {
		for {
			a.UpdateMetrics()

			time.Sleep(a.PollInterval)
		}
	}()
	for {
		time.Sleep(a.ReportInterval)

		a.SendMetrics()
	}
}
