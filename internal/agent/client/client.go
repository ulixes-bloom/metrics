package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	appErrors "github.com/ulixes-bloom/ya-metrics/internal/pkg/errors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/workerpool"
)

type client struct {
	service Service
	http    *http.Client
	conf    *config.Config
}

func New(conf *config.Config, storage service.Storage) *client {
	return &client{
		service: service.New(storage),
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

	// create worker pool for sending metrics to server
	pool := workerpool.New(c.conf.RateLimit, metrics.MetricsCount, c.sendMetric)

	for {
		select {
		case <-reportTicker.C:
			for _, m := range c.service.GetAll() {
				pool.Submit(m)
			}
		case <-ctx.Done():
			log.Debug().Msg("done reporting metrics")
			pool.StopAndWait()
			return
		}
	}
}

func (c *client) sendMetric(m metrics.Metric) error {
	marshalled, err := json.Marshal(m)
	if err != nil {
		return errors.Join(appErrors.ErrFailedMetricMarshall, err)
	}

	buf := bytes.NewBuffer(nil)
	gb := gzip.NewWriter(buf)
	_, err = gb.Write(marshalled)
	if err != nil {
		return errors.Join(appErrors.ErrFailedMetricCompression, err)
	}
	err = gb.Close()
	if err != nil {
		return errors.Join(appErrors.ErrFailedMetricCompression, err)
	}

	url := fmt.Sprintf("%s/update/", c.conf.GetNormilizedServerAddr())
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	req.Header.Add(headers.ContentType, "application/json")
	req.Header.Add(headers.AcceptEncoding, "gzip")
	req.Header.Add(headers.ContentEncoding, "gzip")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code while sending metrics: %s", res.Status)
	}

	defer res.Body.Close()
	return nil
}
