package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
)

func main() {
	conf := parseConfig()

	err := api.Run(conf.RunAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
