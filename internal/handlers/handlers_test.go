package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	type args struct {
		url          string
		method       string
		expectedCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Correct request with gauge metric",
			args: args{
				url:          "/update/gauge/some/1.9",
				method:       http.MethodPost,
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "Correct request with counter metric",
			args: args{
				url:          "/update/counter/some/1",
				method:       http.MethodPost,
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "Empty url",
			args: args{
				url:          "/",
				method:       http.MethodPost,
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Wrong method",
			args: args{
				url:          "/update/gauge/some/1",
				method:       http.MethodGet,
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Wrong metric type",
			args: args{
				url:          "/update/newmetric/some/1",
				method:       http.MethodPost,
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Wrong counter metric value",
			args: args{
				url:          "/update/counter/some/1.89",
				method:       http.MethodPost,
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Url with incomplete metric information",
			args: args{
				url:          "/update/newmetric/1",
				method:       http.MethodPost,
				expectedCode: http.StatusNotFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.args.method, test.args.url, nil)
			w := httptest.NewRecorder()

			UpdateMetric(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.args.expectedCode, res.StatusCode)
		})
	}
}
