package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	RunAddr string `env:"ADDRESS"`
}

func parseConfig() (conf config) {
	flag.StringVar(&conf.RunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()

	env.Parse(&conf)
	return
}
