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
	responseMemory struct {
		status int
		body   bytes.Buffer
	}

	responseWriterWithMemory struct {
		http.ResponseWriter
		responseMemory *responseMemory
	}
)

func newResponseWriterWithMemory(w http.ResponseWriter) *responseWriterWithMemory {
	responseMemory := responseMemory{}
	return &responseWriterWithMemory{
		ResponseWriter: w,
		responseMemory: &responseMemory,
	}
}

// Расширение метода http.ResponseWriter.Write()
func (r *responseWriterWithMemory) Write(b []byte) (int, error) {
	r.responseMemory.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Расширение метода http.ResponseWriter.WriteHeader()
func (r *responseWriterWithMemory) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseMemory.status = statusCode
}

func WithHashing(hashKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Считывание хеша из хедера
			reqHash := r.Header.Get(headers.HashSHA256)
			if reqHash == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Считывание содержимого request body буфера
			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Закрытие request body
			err = r.Body.Close()
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Восстановление содержимого request body буфера
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// Вычисление хеша от тела запроса
			calchash, err := hash.Encode([]byte(reqBody), hashKey)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Сравнение вычисленного и полученного в заголовке хешей
			if reqHash != calchash {
				log.Error().Msg("incorrect hash")
				http.Error(w, "incorrect hash", http.StatusBadRequest)
				return
			}

			wm := newResponseWriterWithMemory(w)

			next.ServeHTTP(wm, r)

			if wm.responseMemory.status == http.StatusOK {
				respBody := wm.responseMemory.body
				respHash, err := hash.Encode(respBody.Bytes(), hashKey)
				if err != nil {
					log.Error().Msg(err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				w.Header().Set(headers.HashSHA256, respHash)
				w.WriteHeader(wm.responseMemory.status)
			}
		})
	}
}
