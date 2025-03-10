package httpapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ulixes-bloom/ya-metrics/internal/server/api"
	"github.com/ulixes-bloom/ya-metrics/internal/server/config"
	"github.com/ulixes-bloom/ya-metrics/internal/server/service"
)

type httpApi struct {
	service api.Service
	conf    *config.Config
	router  *chi.Mux
}

func New(conf *config.Config, storage service.Storage) *httpApi {
	srv := service.New(storage, conf)
	newAPI := httpApi{
		service: srv,
		conf:    conf,
	}
	newAPI.router = newAPI.newRouter()
	return &newAPI
}

func (a *httpApi) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- http.ListenAndServe(a.conf.RunAddr, a.router)
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("httpapi.run: %w", err)
	case <-ctx.Done():
		return a.service.Shutdown(ctx)
	}
}
