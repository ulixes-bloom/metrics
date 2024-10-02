package main

import (
	"net/http"

	"github.com/ulixes-bloom/ya-metrics/internal/server"
)

func main() {
	r := server.Router()
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
