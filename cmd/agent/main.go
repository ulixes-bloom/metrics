package main

import (
	"log"

	"github.com/ulixes-bloom/ya-metrics/internal/agent/client"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
)

func main() {
	conf, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	client.New(conf).Run()
}
