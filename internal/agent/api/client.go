package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type Client struct {
	Service        Service
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerAddr     string
}

func NewClient(service Service, pollInterval, reportInterval time.Duration, serverAddr string) *Client {
	return &Client{
		Service:        service,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		ServerAddr:     serverAddr,
	}
}

func (c *Client) UpdateMetrics() {
	c.Service.UpdateMetrics()
}

func (c *Client) SendMetrics() {
	for _, v := range c.Service.GetAll() {
		c.SendMetric(v)
	}
}

func (c *Client) SendMetric(m metrics.Metric) {
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("impossible to marshall metric: %s", err)
	}

	buf := bytes.NewBuffer(nil)
	gb := gzip.NewWriter(buf)
	_, err = gb.Write(marshalled)
	if err != nil {
		log.Fatalf("impossible to compress metric ussing gzip: %s", err)
	}
	err = gb.Close()
	if err != nil {
		log.Fatalf("impossible to compress metric ussing gzip: %s", err)
	}

	url := fmt.Sprintf("%s/update/", c.ServerAddr)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Encoding", "gzip")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
}
