package main

import (
	"flag"

	"github.com/caarlos0/env"
)

type config struct {
	ServerAddr     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func parseConfig() (conf config) {
	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "address and port of server")
	flag.IntVar(&conf.ReportInterval, "r", 10, "metrics report interval")
	flag.IntVar(&conf.PollInterval, "p", 2, "metrics update interval")
	flag.Parse()

	env.Parse(&conf)
	return
}
