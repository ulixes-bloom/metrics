package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/memory"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metricerrors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type client struct {
	service Service
	conf    *config.Config
}

func New(conf *config.Config) *client {
	ms := memory.NewStorage()
	s := service.New(ms)

	return &client{
		service: s,
		conf:    conf,
	}
}

func (c *client) Run(ctx context.Context) {
	pollTicker := time.NewTicker(c.conf.GetPollIntervalDuration())
	reportTicker := time.NewTicker(c.conf.GetReportIntervalDuration())
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			c.UpdateMetrics()
		case <-reportTicker.C:
			if err := c.SendMetrics(); err != nil {
				log.Error().Msg(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *client) UpdateMetrics() {
	c.service.UpdateMetrics()
}

func (c *client) SendMetrics() error {
	metricsSlice := []metrics.Metric{}
	for _, v := range c.service.GetAll() {
		metricsSlice = append(metricsSlice, v)
	}

	marshalled, err := json.Marshal(metricsSlice)
	if err != nil {
		return errors.Join(metricerrors.ErrFailedMetricMarshall, err)
	}

	buf := bytes.NewBuffer(nil)
	gb := gzip.NewWriter(buf)
	_, err = gb.Write(marshalled)
	if err != nil {
		return errors.Join(metricerrors.ErrFailedMetricCompression, err)
	}
	err = gb.Close()
	if err != nil {
		return errors.Join(metricerrors.ErrFailedMetricCompression, err)
	}

	url := fmt.Sprintf("%s/updates/", c.conf.GetNormilizedServerAddr())
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	req.Header.Add(headers.ContentType, "application/json")
	req.Header.Add(headers.AcceptEncoding, "gzip")
	req.Header.Add(headers.ContentEncoding, "gzip")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code while sending metrics: %s", res.Status)
	}

	defer res.Body.Close()
	return nil
}
