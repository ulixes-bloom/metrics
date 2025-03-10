package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/hash"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
)

type (
	// responseMemory stores the status and body of an HTTP response to capture
	// the response data for further processing (e.g., hashing).
	responseMemory struct {
		status int
		body   bytes.Buffer
	}

	// responseWriterWithMemory wraps an http.ResponseWriter and includes a
	// responseMemory to capture the response body and status for later use.
	responseWriterWithMemory struct {
		http.ResponseWriter
		responseMemory *responseMemory
	}
)

func newResponseWriterWithMemory(w http.ResponseWriter) *responseWriterWithMemory {
	respMemory := responseMemory{}
	return &responseWriterWithMemory{
		ResponseWriter: w,
		responseMemory: &respMemory,
	}
}

func (r *responseWriterWithMemory) Write(b []byte) (int, error) {
	r.responseMemory.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseWriterWithMemory) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseMemory.status = statusCode
}

// WithHashing is a middleware that validates the hash of the request body
// and appends the hash of the response body to the response header. It uses
// the provided hash key to calculate and compare hashes.
func WithHashing(hashKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the hash value from the requst header
			reqHash := r.Header.Get(headers.HashSHA256)
			if reqHash == "" {
				next.ServeHTTP(w, r)
				return
			}

			// read the request body into a buffer
			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			// restore the request body so it can be read by subsequent handlers
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// calculate the hash value of the request body using the provided hash key
			calchash, err := hash.Encode(reqBody, hashKey)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// compare the calculated hash with the hash provided in the header
			if reqHash != calchash {
				log.Error().Msg("incorrect hash")
				http.Error(w, "incorrect hash", http.StatusBadRequest)
				return
			}

			// create a wrapper around the ResponseWriter to capture the response body and status
			wm := newResponseWriterWithMemory(w)

			next.ServeHTTP(wm, r)
			// if the response status is HTTP 200 (OK), calculate the response body hash
			if wm.responseMemory.status == http.StatusOK {
				respBody := wm.responseMemory.body
				respHash, err := hash.Encode(respBody.Bytes(), hashKey)
				if err != nil {
					log.Error().Msg(err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				w.Header().Set(headers.HashSHA256, respHash)
			}
		})
	}
}
