package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
)

type api struct {
	service Service
	conf    *config.Config
	router  *chi.Mux
}

func New(conf *config.Config, storage service.Storage) *api {
	srv := service.New(storage, conf)
	newAPI := api{
		service: srv,
		conf:    conf,
	}
	newAPI.router = newAPI.newRouter()
	return &newAPI
}

func (a *api) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- http.ListenAndServe(a.conf.RunAddr, a.router)
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("api.run: %w", err)
	case <-ctx.Done():
		return a.service.Shutdown(ctx)
	}
}
