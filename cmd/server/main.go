package main

import (
	"log"

	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
)

func main() {
	conf := parseConfig()

	err := api.Run(conf.RunAddr, conf.LogLvl)
	if err != nil {
		log.Fatal(err)
	}
}
