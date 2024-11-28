package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metricerrors"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

// worker принимает на вход канал с входными данными для отправки метрик
func (c *client) worker(metricsToSend <-chan metrics.Metric) {
	for {
		m := <-metricsToSend
		err := c.sendMetric(m)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func (c *client) sendMetric(m metrics.Metric) error {
	marshalled, err := json.Marshal(m)
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
