package client

import (
	"context"
	"net/http"
	"sync"

	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/memory"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type client struct {
	service Service
	http    *http.Client
	conf    *config.Config
}

func New(conf *config.Config) *client {
	ms := memory.NewStorage()

	return &client{
		service: service.New(ms),
		http:    &http.Client{},
		conf:    conf,
	}
}

func (c *client) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		c.pollMetrics(ctx)
		wg.Done()
	}()

	go func() {
		c.reportMetrics(ctx)
		wg.Done()
	}()

	wg.Wait()
}

func (c *client) pollMetrics(ctx context.Context) {
	pollTicker := time.NewTicker(c.conf.GetPollIntervalDuration())
	defer pollTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			err := c.service.Poll(ctx)
			if err != nil {
				log.Error().Msg(err.Error())
			}
		case <-ctx.Done():
			log.Debug().Msg("done polling metrics")
			return
		}
	}
}

func (c *client) reportMetrics(ctx context.Context) {
	reportTicker := time.NewTicker(c.conf.GetReportIntervalDuration())
	defer reportTicker.Stop()

	// создаем канал, через который будем отправлять метрики в worker'ы
	metricsToSendChan := make(chan metrics.Metric, metrics.MetricsCount)
	defer close(metricsToSendChan)

	// создаем worker'ов для отправки метрик на сервер
	for w := 1; w <= c.conf.RateLimit; w++ {
		go c.worker(metricsToSendChan)
	}

	for {
		select {
		case <-reportTicker.C:
			for _, m := range c.service.GetAll() {
				metricsToSendChan <- m
			}
		case <-ctx.Done():
			log.Debug().Msg("done reporting metrics")
			return
		}
	}
}
