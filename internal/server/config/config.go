// Package config provides functionality for parsing server configuration
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

type Config struct {
	RunAddr         string `env:"ADDRESS" json:"address"`                     // The address and port for the server to listen on.
	LogLvl          string `env:"LOGLVL" json:"loglvl"`                       // The logging level to be used (e.g., Info, Debug).
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`       // Interval at which metrics are stored.
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"` // Path to store metrics data in a file.
	Restore         bool   `env:"RESTORE" json:"restore"`                     // Flag to determine if metrics should be restored from storage.
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`           // Data source name for connecting to a PostgreSQL database.
	HashKey         string `env:"KEY" json:"hash_key"`                        // Key used for signing and validating metrics data.
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`               // Public key for data encryption.
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

	flag.StringVar(&conf.RunAddr, "a", conf.RunAddr, "address and port to run server")
	flag.StringVar(&conf.LogLvl, "l", conf.LogLvl, "logging level")
	flag.IntVar(&conf.StoreInterval, "i", conf.StoreInterval, "store interval")
	flag.StringVar(&conf.FileStoragePath, "f", conf.FileStoragePath, "file storage path")
	flag.BoolVar(&conf.Restore, "r", conf.Restore, "to restore metrics data")
	flag.StringVar(&conf.DatabaseDSN, "d", conf.DatabaseDSN, "Postgresql data source name")
	flag.StringVar(&conf.HashKey, "k", conf.HashKey, "key to sign the metrics data")
	flag.StringVar(&conf.CryptoKey, "crypto-key", conf.HashKey, "public key for data encryption")
	flag.StringVar(&configFile, "c", configFile, "json file with configuration")
	flag.Parse()

	err = env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}

	if conf.StoreInterval < 0 {
		return nil, errors.New("config.parse: negative store interval")
	}

	fmt.Println(conf)
	return &conf, nil
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

// GetDefault returns a Config object populated with default values for all
// configuration options. These defaults are used when no other values are
// provided by the user through command-line flags or environment variables.
func GetDefault() (conf *Config) {
	return &Config{
		RunAddr:         ":8080",
		LogLvl:          "Info",
		StoreInterval:   300,
		FileStoragePath: "metrics_store.txt",
		Restore:         true,
		DatabaseDSN:     "",
		HashKey:         "",
		CryptoKey:       "",
	}
}

// GetStoreIntervalDuration converts the StoreInterval field to a time.Duration.
func (c *Config) GetStoreIntervalDuration() time.Duration {
	return time.Duration(c.StoreInterval) * time.Second
}
