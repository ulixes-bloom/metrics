package api

import (
	"net/http"
	"time"
)

func (handler *Handler) WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		duration := time.Since(start)
		handler.Logger.Info().
			Str("uri", uri).
			Str("method", method).
			Str("duration", duration.String()).
			Msg("got incoming HTTP request")
	})
}
