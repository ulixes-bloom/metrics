// Package config provides functionality for parsing client configuration
// from command-line flags and environment variables.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/caarlos0/env"
)

// Config holds the configuration parameters for the client.
// Parameters can be set through environment variables or command-line flags.
type Config struct {
	ServerAddr     string `env:"ADDRESS" json:"address"`                 // Address and port of the server.
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"` // Interval for reporting metrics, in seconds.
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`     // Interval for polling metrics, in seconds.
	RateLimit      int    `env:"RATE_LIMIT" json:"rate_limit"`           // Rate limit for metric updates.
	LogLvl         string `env:"LOGLVL" json:"loglvl"`                   // Logging level (e.g., "info", "debug").
	HashKey        string `env:"KEY" json:"key"`                         // Key for signing metrics data.
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`           // Public key for data encryption.
}

// Parse parses the configuration from command-line flags and environment variables.
func Parse() (*Config, error) {
	var conf Config
	var err error

	// read config from json config file
	configFile := getConfigFileName()
	if configFile != "" {
		conf, err = parseConfigFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("config.parse: %w", err)
		}
	}

	// fill config empty parameters with default values
	defaultValues := GetDefault()
	err = mergo.Merge(&conf, defaultValues)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}

	flag.StringVar(&conf.ServerAddr, "a", conf.ServerAddr, "address and port of server")
	flag.IntVar(&conf.ReportInterval, "r", conf.ReportInterval, "metrics report interval")
	flag.IntVar(&conf.PollInterval, "p", conf.PollInterval, "metrics update interval")
	flag.IntVar(&conf.RateLimit, "l", conf.RateLimit, "rate limit for metric updates")
	flag.StringVar(&conf.LogLvl, "ll", conf.LogLvl, "logging level")
	flag.StringVar(&conf.HashKey, "k", conf.HashKey, "key to sign the metrics data")
	flag.StringVar(&conf.CryptoKey, "crypto-key", conf.HashKey, "public key for data encryption")
	flag.StringVar(&configFile, "c", configFile, "json file with configuration")
	flag.Parse()

	err = env.Parse(&conf)
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

	fmt.Println(conf)
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

func getConfigFileName() (configPath string) {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-c=") {
			configPath = strings.TrimPrefix(arg, "-c=")
		}
		if strings.HasPrefix(arg, "-config=") {
			configPath = strings.TrimPrefix(arg, "-config=")
		}
	}
	if env, isExist := os.LookupEnv("CONFIG"); isExist {
		configPath = env
	}

	return
}

func parseConfigFromFile(fname string) (conf Config, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return conf, fmt.Errorf("config.parseConfigFromFile: %w", err)
	}

	dec := json.NewDecoder(f)
	if err := dec.Decode(&conf); err != nil {
		return conf, fmt.Errorf("config.parseConfigFromFile: %w", err)
	}

	return conf, nil
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
