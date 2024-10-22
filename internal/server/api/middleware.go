package api

import (
	"net/http"
	"time"
)

func (h *Handler) WithLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		handler.ServeHTTP(w, r)

		duration := time.Since(start)
		h.Logger.Debug().
			Str("uri", uri).
			Str("method", method).
			Str("duration", duration.String()).
			Msg("got incoming HTTP request")
	})
}
