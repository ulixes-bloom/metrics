package api

import (
	"bytes"
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

	url := fmt.Sprintf("%s/update/", c.ServerAddr)
	resp, err := http.Post(url, "application/json", bytes.NewReader(marshalled))
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
