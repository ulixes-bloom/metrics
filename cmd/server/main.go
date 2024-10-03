package main

import (
	"net/http"

	"github.com/ulixes-bloom/ya-metrics/internal/server"
)

func main() {
	parseFlags()

	r := server.Router()
	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
