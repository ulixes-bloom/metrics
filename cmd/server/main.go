package main

import (
	"net/http"

	"github.com/ulixes-bloom/ya-metrics/internal/server"
)

func main() {
	conf := parseConfig()

	r := server.Router()
	err := http.ListenAndServe(conf.runAddr, r)
	if err != nil {
		panic(err)
	}
}
