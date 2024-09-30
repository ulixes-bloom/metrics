package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/ulixes-bloom/ya-metrics/internal/storage"
)

var MemStorage = storage.NewMemStorage()

func UpdateMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		url := req.URL.Path
		if pattern := regexp.MustCompile(`^/update\/(.*?)\/(.*?)\/(.*)$`); pattern.MatchString(url) {
			match := pattern.FindStringSubmatch(url)
			mtype, mname, mval := match[1], match[2], match[3]

			switch mtype {
			case "gauge":
				if val, err := strconv.ParseFloat(mval, 64); err == nil {
					MemStorage.AddGauge(mname, val)
				} else {
					res.WriteHeader(http.StatusBadRequest)
				}
			case "counter":
				if val, err := strconv.ParseInt(mval, 10, 64); err == nil {
					MemStorage.AddCounter(mname, val)
				} else {
					res.WriteHeader(http.StatusBadRequest)
				}
			default:
				res.WriteHeader(http.StatusBadRequest)
			}

			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusNotFound)
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
