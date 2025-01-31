// Package config provides functionality for parsing server configuration
// from command-line flags and environment variables.
package config

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddr         string        `env:"ADDRESS"`           // The address and port for the server to listen on.
	LogLvl          string        `env:"LOGLVL"`            // The logging level to be used (e.g., Info, Debug).
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`    // Interval at which metrics are stored.
	FileStoragePath string        `env:"FILE_STORAGE_PATH"` // Path to store metrics data in a file.
	Restore         bool          `env:"RESTORE"`           // Flag to determine if metrics should be restored from storage.
	DatabaseDSN     string        `env:"DATABASE_DSN"`      // Data source name for connecting to a PostgreSQL database.
	HashKey         string        `env:"KEY"`               // Key used for signing and validating metrics data.
}

// Parse parses the configuration from command-line flags and environment variables.
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

// GetDefault returns a Config object populated with default values for all
// configuration options. These defaults are used when no other values are
// provided by the user through command-line flags or environment variables.
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
