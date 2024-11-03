package compress

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
)

type GzipWriter struct {
	w      http.ResponseWriter
	Writer *gzip.Writer
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:      w,
		Writer: gzip.NewWriter(w),
	}
}

func (gw *GzipWriter) Header() http.Header {
	return gw.w.Header()
}

func (gw *GzipWriter) Write(p []byte) (int, error) {
	return gw.Writer.Write(p)
}

func (gw *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		gw.w.Header().Add(headers.ContentEncoding, "gzip")
	}
	gw.w.WriteHeader(statusCode)
}

func (gw *GzipWriter) Close() error {
	return gw.Writer.Close()
}

type GzipReader struct {
	r      io.ReadCloser
	Reader *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:      r,
		Reader: zr,
	}, nil
}

func (gr GzipReader) Read(p []byte) (n int, err error) {
	return gr.Reader.Read(p)
}

func (gr *GzipReader) Close() error {
	if err := gr.r.Close(); err != nil {
		return err
	}
	return gr.Reader.Close()
}
