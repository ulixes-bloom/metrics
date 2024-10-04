package main

import (
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/agent"
)

func main() {
	conf := parseConfig()

	a := agent.NewAgent(
		time.Duration(conf.pollInterval)*time.Second,
		time.Duration(conf.reportInterval)*time.Second,
		"http://"+conf.serverAddr)

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
