package config

import (
	"errors"
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLvl         string `env:"LOGLVL"`
	HashKey        string `env:"KEY"`
}

func Parse() (*Config, error) {
	var conf Config

	flag.StringVar(&conf.ServerAddr, "a", "localhost:8080", "address and port of server")
	flag.IntVar(&conf.ReportInterval, "r", 10, "metrics report interval")
	flag.IntVar(&conf.PollInterval, "p", 2, "metrics update interval")
	flag.StringVar(&conf.LogLvl, "l", "debug", "logging level")
	flag.StringVar(&conf.HashKey, "k", "", "key to sign the metrics data")
	flag.Parse()

	env.Parse(&conf)

	if conf.ReportInterval <= 0 {
		return nil, errors.New("negative report interval")
	}
	if conf.PollInterval <= 0 {
		return nil, errors.New("negative poll interval")
	}

	return &conf, nil
}

func (c *Config) GetNormilizedServerAddr() string {
	return "http://" + c.ServerAddr
}

func (c *Config) GetReportIntervalDuration() time.Duration {
	return time.Duration(c.ReportInterval) * time.Second
}

func (c *Config) GetPollIntervalDuration() time.Duration {
	return time.Duration(c.PollInterval) * time.Second
}
