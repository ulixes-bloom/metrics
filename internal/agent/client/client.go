package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/agent/config"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/memory"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

type client struct {
	Service        Service
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerAddr     string
}

func New(conf config.Config) *client {
	ms := memory.NewStorage()
	s := service.New(ms)

	return &client{
		Service:        s,
		PollInterval:   time.Duration(conf.PollInterval) * time.Second,
		ReportInterval: time.Duration(conf.ReportInterval) * time.Second,
		ServerAddr:     "http://" + conf.ServerAddr,
	}
}

func (c *client) Run() {
	go func() {
		for {
			c.UpdateMetrics()

			time.Sleep(c.PollInterval)
		}
	}()
	for {
		time.Sleep(c.ReportInterval)

		c.SendMetrics()
	}
}

func (c *client) UpdateMetrics() {
	c.Service.UpdateMetrics()
}

func (c *client) SendMetrics() {
	for _, v := range c.Service.GetAll() {
		c.SendMetric(v)
	}
}

func (c *client) SendMetric(m metrics.Metric) {
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
	req.Header.Add(headers.ContentType, "application/json")
	req.Header.Add(headers.AcceptEncoding, "gzip")
	req.Header.Add(headers.ContentEncoding, "gzip")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
}
