package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ulixes-bloom/ya-metrics/internal/agent/client"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
)

func main() {
	conf, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	
	client.New(conf).Run(ctx)
}
