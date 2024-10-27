package main

import (
	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
)

func main() {
	conf := config.Parse()
	api.New(conf).Run()
}
