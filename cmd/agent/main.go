package main

import (
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/agent"
)

func main() {
	parseFlags()

	a := agent.NewAgent(
		time.Duration(flagPollInterval)*time.Second,
		time.Duration(flagReportInterval)*time.Second,
		"http://"+flagServerAddr)

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
