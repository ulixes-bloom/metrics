package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/compress"
)

func (h *Handler) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		h.Logger.Debug().
			Str("uri", uri).
			Str("method", method).
			Str("duration", duration.String()).
			Msg("got incoming HTTP request")
	})
}

func (h *Handler) WithCompressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(acceptEncoding, "gzip") &&
			(strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")) {
			cw := compress.NewGzipWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := compress.NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
