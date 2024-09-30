package main

import (
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/agent"
)

func main() {
	a := agent.NewAgent(
		time.Duration(2*time.Second),
		time.Duration(10*time.Second),
		"http://localhost:8080")

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
