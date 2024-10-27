package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddr         string        `env:"ADDRESS"`
	LogLvl          string        `env:"LOGLVL"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
}

func Parse() (conf Config) {
	flag.StringVar(&conf.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&conf.LogLvl, "l", "Info", "logging level")
	flag.DurationVar(&conf.StoreInterval, "i", 300*time.Second, "store interval")
	flag.StringVar(&conf.FileStoragePath, "f", "metrics_store.txt", "file storage path")
	flag.BoolVar(&conf.Restore, "r", true, "to restore metrics data")
	flag.Parse()

	env.Parse(&conf)
	return
}
