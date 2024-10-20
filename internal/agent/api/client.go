package api

import (
	"fmt"
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
	for k, v := range c.Service.GetAllGauges() {
		c.SendMetric(metrics.Gauge, k, fmt.Sprintf("%v", v))
	}

	for k, v := range c.Service.GetAllCounters() {
		c.SendMetric(metrics.Counter, k, fmt.Sprintf("%d", v))
	}
}

func (c *Client) SendMetric(mtype, mname, mval string) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", c.ServerAddr, mtype, mname, mval)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
}
