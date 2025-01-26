package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// WithLogging is a middleware that logs details about incoming HTTP requests.
// It logs the request URI, HTTP method, and the duration it took to process the request.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Debug().
			Str("uri", uri).
			Str("method", method).
			Str("duration", duration.String()).
			Msg("got incoming HTTP request")
	})
}
