package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
)

func main() {
	conf := config.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	api.New(conf).Run(ctx)
}
