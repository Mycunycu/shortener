package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "happy path",
			request: "localhost:8080",
			body:    "https://github.com/",
			want: want{
				contentType: "text/html; charset=UTF-8",
				statusCode:  201,
				body:        "http://localhost:8080/1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := ShortenURL()

			h.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(bodyResult))
		})
	}
}

func TestExpandURL(t *testing.T) {
	type want struct {
		statusCode  int
		headerValue string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "happy path",
			request: "localhost:8080",
			want: want{
				statusCode:  307,
				headerValue: "https://github.com/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "0")

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h := ExpandURL()

			h.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.headerValue, result.Header.Get("Location"))
		})
	}
}
