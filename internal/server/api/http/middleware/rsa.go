package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/rsa"
)

// WithRSA is a middleware that decrypts request body.
// It uses the provided cryptoKey.
func WithRSA(cryptoKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the request body into a buffer
			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			// decrypt body
			cleartext, err := rsa.Decrypt([]byte(reqBody), cryptoKey)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// restore decrypted request body so it can be read by subsequent handlers
			r.Body = io.NopCloser(bytes.NewBuffer(cleartext))

			next.ServeHTTP(w, r)
		})
	}
}
