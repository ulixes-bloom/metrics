package main

import (
	"net/http"

	"github.com/ulixes-bloom/ya-metrics/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.UpdateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
