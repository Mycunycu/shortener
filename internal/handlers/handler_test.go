package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mycunycu/shortener/internal/config"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg = config.Config{
	ServerAddress:   "localhost:8080",
	BaseURL:         "http://localhost:8080",
	FileStoragePath: "./storage.txt",
}

func TestShortenURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name string
		path string
		body string
		want want
	}{
		{
			name: "happy path",
			path: cfg.ServerAddress,
			body: "https://test.com",
			want: want{
				contentType: "text/html; charset=UTF-8",
				statusCode:  201,
				body:        "http://localhost:8080/1",
			},
		},
		{
			name: "empty body",
			path: cfg.ServerAddress,
			body: "",
			want: want{
				contentType: "",
				statusCode:  400,
				body:        "",
			},
		},
		{
			name: "invalid body",
			path: cfg.ServerAddress,
			body: "https:/test.com",
			want: want{
				contentType: "",
				statusCode:  400,
				body:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			storage, _ := repository.NewStorage(cfg.FileStoragePath)
			repo := repository.NewShortURL(storage)

			h := NewHandler(cfg, repo, storage).ShortenURL()

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
	path := "localhost:8080"

	type want struct {
		statusCode  int
		headerValue string
	}

	tests := []struct {
		name string
		path string
		id   string
		want want
	}{
		{
			name: "happy path",
			path: path,
			id:   "1",
			want: want{
				statusCode:  307,
				headerValue: "https://test.com",
			},
		},
		{
			name: "no have data",
			path: path,
			id:   "2",
			want: want{
				statusCode:  204,
				headerValue: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			storage, _ := repository.NewStorage(cfg.FileStoragePath)
			repo := repository.NewShortURL(storage)
			repo.Set("https://test.com")

			h := NewHandler(cfg, repo, storage).ExpandURL()
			h.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.headerValue, result.Header.Get("Location"))
		})
	}
}

func TestShorten(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	path := "localhost:8080/api/shorten"

	tests := []struct {
		name string
		path string
		body string
		want want
	}{
		{
			name: "happy path",
			path: path,
			body: `{"url": "http://test.ru"}`,
			want: want{
				contentType: "application/json",
				statusCode:  201,
				body:        `{"result": "http://localhost:8080/1"}`,
			},
		},
		// TODO add test cases
		// {
		// 	name: "empty body",
		// 	path: path,
		// 	body: "",
		// 	want: want{
		// 		contentType: "",
		// 		statusCode:  400,
		// 		body:        "",
		// 	},
		// },
		// {
		// 	name: "invalid body",
		// 	path: path,
		// 	body: "https:/test.com",
		// 	want: want{
		// 		contentType: "",
		// 		statusCode:  400,
		// 		body:        "",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			storage, _ := repository.NewStorage(cfg.FileStoragePath)
			repo := repository.NewShortURL(storage)
			h := NewHandler(cfg, repo, storage).Shorten()
			h.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.JSONEq(t, tt.want.body, string(bodyResult))
		})
	}
}
