package main

import (
	"flag"
)

var flagServerAddr string
var flagReportInterval uint
var flagPollInterval uint

func parseFlags() {
	flag.StringVar(&flagServerAddr, "a", "http://localhost:8080", "address and port of server")
	flag.UintVar(&flagReportInterval, "r", 10, "metric report interval")
	flag.UintVar(&flagPollInterval, "p", 2, "metrics update interval")

	flag.Parse()
}
