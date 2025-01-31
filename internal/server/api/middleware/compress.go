package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
)

type (
	// gzipWriter wraps an http.ResponseWriter and a gzip.Writer to handle
	// gzipped response content.
	gzipWriter struct {
		w      http.ResponseWriter
		Writer *gzip.Writer
	}

	// gzipReader wraps an io.ReadCloser and a gzip.Reader to handle
	// reading gzipped request bodies.
	gzipReader struct {
		r      io.ReadCloser
		Reader *gzip.Reader
	}
)

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		w:      w,
		Writer: gzip.NewWriter(w),
	}
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("middleware.newGzipReader: %w", err)
	}
	return &gzipReader{
		r:      r,
		Reader: zr,
	}, nil
}

func (gw *gzipWriter) Header() http.Header {
	return gw.w.Header()
}

func (gw *gzipWriter) Write(p []byte) (int, error) {
	return gw.Writer.Write(p)
}

// WriteHeader sets the HTTP status code and, if the status code is below 300,
// adds the "Content-Encoding: gzip" header to indicate that the response is gzipped.
func (gw *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		gw.w.Header().Set(headers.ContentEncoding, "gzip")
	}
	gw.w.WriteHeader(statusCode)
}

func (gw *gzipWriter) Close() error {
	return gw.Writer.Close()
}

func (gr gzipReader) Read(p []byte) (n int, err error) {
	return gr.Reader.Read(p)
}

func (gr *gzipReader) Close() error {
	if err := gr.r.Close(); err != nil {
		return fmt.Errorf("middleware.gzipReader.close: %w", err)
	}
	return gr.Reader.Close()
}

// WithCompressing is a middleware that compresses HTTP responses with gzip
// if the client supports it, and decompresses gzipped request bodies if the
// client sends them. It adds the appropriate "Content-Encoding" and "Accept-Encoding"
// headers to indicate compression.
func WithCompressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get(headers.AcceptEncoding)
		if strings.Contains(acceptEncoding, "gzip") {
			cw := newGzipWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get(headers.ContentEncoding)
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newGzipReader(r.Body)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
