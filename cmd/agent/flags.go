package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	serverAddr     string `env:"ADDRESS"`
	reportInterval uint   `env:"REPORT_INTERVAL"`
	pollInterval   uint   `env:"POLL_INTERVAL"`
}

func parseConfig() (conf config) {
	flag.StringVar(&conf.serverAddr, "a", "localhost:8080", "address and port of server")
	flag.UintVar(&conf.reportInterval, "r", 10, "metrics report interval")
	flag.UintVar(&conf.pollInterval, "p", 2, "metrics update interval")
	flag.Parse()

	env.Parse(&conf)
	return
}
