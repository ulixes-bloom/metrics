// Package config provides functionality for parsing client configuration
// from command-line flags and environment variables.
package config

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

// Config holds the configuration parameters for the client.
// Parameters can be set through environment variables or command-line flags.
type Config struct {
	ServerAddr     string `env:"ADDRESS"`         // Address and port of the server.
	ReportInterval int    `env:"REPORT_INTERVAL"` // Interval for reporting metrics, in seconds.
	PollInterval   int    `env:"POLL_INTERVAL"`   // Interval for polling metrics, in seconds.
	RateLimit      int    `env:"RATE_LIMIT"`      // Rate limit for metric updates.
	LogLvl         string `env:"LOGLVL"`          // Logging level (e.g., "info", "debug").
	HashKey        string `env:"KEY"`             // Key for signing metrics data.
	CryptoKey      string `env:"CRYPTO_KEY"`      // Public key for data encryption.
}

// Parse parses the configuration from command-line flags and environment variables.
func Parse() (*Config, error) {
	var conf Config
	defaultValues := GetDefault()

	flag.StringVar(&conf.ServerAddr, "a", defaultValues.ServerAddr, "address and port of server")
	flag.IntVar(&conf.ReportInterval, "r", defaultValues.ReportInterval, "metrics report interval")
	flag.IntVar(&conf.PollInterval, "p", defaultValues.PollInterval, "metrics update interval")
	flag.IntVar(&conf.RateLimit, "l", defaultValues.RateLimit, "rate limit for metric updates")
	flag.StringVar(&conf.LogLvl, "ll", defaultValues.LogLvl, "logging level")
	flag.StringVar(&conf.HashKey, "k", defaultValues.HashKey, "key to sign the metrics data")
	flag.StringVar(&conf.CryptoKey, "crypto-key", defaultValues.HashKey, "public key for data encryption")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}

	if conf.ReportInterval <= 0 {
		return nil, errors.New("config.parse: negative or zero report interval")
	}
	if conf.PollInterval <= 0 {
		return nil, errors.New("config.parse: negative or zero poll interval")
	}
	if conf.RateLimit <= 0 {
		return nil, errors.New("config.parse: negative or zero rate interval")
	}

	return &conf, nil
}

// GetDefault returns a Config object populated with default values for all
// configuration options. These defaults are used when no other values are
// provided by the user through command-line flags or environment variables.
func GetDefault() (conf *Config) {
	return &Config{
		ServerAddr:     "localhost:8080",
		ReportInterval: 10,
		PollInterval:   2,
		RateLimit:      1,
		LogLvl:         "info",
		HashKey:        "",
		CryptoKey:      "",
	}
}

// GetNormilizedServerAddr returns the server address normalized with an "http://" prefix.
func (c *Config) GetNormilizedServerAddr() string {
	return "http://" + c.ServerAddr
}

// GetReportIntervalDuration converts the ReportInterval field to a time.Duration.
func (c *Config) GetReportIntervalDuration() time.Duration {
	return time.Duration(c.ReportInterval) * time.Second
}

// GetPollIntervalDuration converts the PollInterval field to a time.Duration.
func (c *Config) GetPollIntervalDuration() time.Duration {
	return time.Duration(c.PollInterval) * time.Second
}
