package main

import (
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/agent"
)

func main() {
	conf := parseConfig()

	a := agent.NewAgent(
		time.Duration(conf.PollInterval)*time.Second,
		time.Duration(conf.ReportInterval)*time.Second,
		"http://"+conf.ServerAddr)

	go func() {
		for {
			a.UpdateGuageMetrics()
			a.UpdateCounterMetrics()

			time.Sleep(a.PollInterval)
		}
	}()
	for {
		time.Sleep(a.ReportInterval)

		a.SendMetrics()
	}
}
