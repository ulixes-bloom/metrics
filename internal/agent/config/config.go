package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func Parse() (conf Config, err error) {
	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "address and port of server")
	flag.IntVar(&conf.ReportInterval, "r", 10, "metrics report interval")
	flag.IntVar(&conf.PollInterval, "p", 2, "metrics update interval")
	flag.Parse()

	env.Parse(&conf)

	if conf.ReportInterval <= 0 {
		err = errors.New("negative report interval")
	}
	if conf.PollInterval <= 0 {
		err = errors.New("negative poll interval")
	}

	return
}
