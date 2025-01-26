package config

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddr         string        `env:"ADDRESS"`
	LogLvl          string        `env:"LOGLVL"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	DatabaseDSN     string        `env:"DATABASE_DSN"`
	HashKey         string        `env:"KEY"`
}

func Parse() (*Config, error) {
	var conf Config
	defaultValues := GetDefault()

	flag.StringVar(&conf.RunAddr, "a", defaultValues.RunAddr, "address and port to run server")
	flag.StringVar(&conf.LogLvl, "l", defaultValues.LogLvl, "logging level")
	flag.DurationVar(&conf.StoreInterval, "i", defaultValues.StoreInterval, "store interval")
	flag.StringVar(&conf.FileStoragePath, "f", defaultValues.FileStoragePath, "file storage path")
	flag.BoolVar(&conf.Restore, "r", defaultValues.Restore, "to restore metrics data")
	flag.StringVar(&conf.DatabaseDSN, "d", defaultValues.DatabaseDSN, "Postgresql data source name")
	flag.StringVar(&conf.HashKey, "k", defaultValues.HashKey, "key to sign the metrics data")

	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}

	if conf.StoreInterval <= 0 {
		return nil, errors.New("config.parse: negative or zero store interval")
	}

	return &conf, nil
}

func GetDefault() (conf *Config) {
	return &Config{
		RunAddr:         ":8080",
		LogLvl:          "Info",
		StoreInterval:   300 * time.Second,
		FileStoragePath: "metrics_store.txt",
		Restore:         true,
		DatabaseDSN:     "",
		HashKey:         "",
	}
}
