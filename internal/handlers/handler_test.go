package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Mycunycu/shortener/internal/config"
	"github.com/Mycunycu/shortener/internal/models"
	"github.com/Mycunycu/shortener/internal/repository/mocks"
	"github.com/Mycunycu/shortener/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var cfg = config.Config{
	ServerAddress:   "localhost:8080",
	BaseURL:         "http://localhost:8080",
	FileStoragePath: "./storage.txt",
	DatabaseDSN:     "postgres://user:123@localhost:5432/practicum",
}
var timeout = time.Duration(time.Second * 5)

func TestShortenURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
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
			},
		},
		{
			name: "empty body",
			path: cfg.ServerAddress,
			body: "",
			want: want{
				contentType: "",
				statusCode:  400,
			},
		},
		{
			name: "invalid body",
			path: cfg.ServerAddress,
			body: "https:/test.com",
			want: want{
				contentType: "",
				statusCode:  400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			repo := &mocks.Repositorier{}
			shortURL := services.NewShortURL(cfg.BaseURL, repo)
			h := NewHandler(shortURL, timeout).ShortenURL()

			repo.On("Save", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.ShortenEty")).Return(nil)

			h.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
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
		// {
		// 	name: "no have data",
		// 	path: path,
		// 	id:   "a",
		// 	want: want{
		// 		statusCode:  204,
		// 		headerValue: "",
		// 	},
		// },
		// {
		// 	name: "no have id",
		// 	path: path,
		// 	id:   "",
		// 	want: want{
		// 		statusCode:  400,
		// 		headerValue: "",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			repo := &mocks.Repositorier{}
			shortURL := services.NewShortURL(cfg.BaseURL, repo)
			h := NewHandler(shortURL, timeout).ExpandURL()

			repo.On("GetByShortID", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string")).Return(
				models.ShortenEty{
					ShortID:     "1",
					OriginalURL: "https://test.com",
				}, nil)
			// repo.On("GetByShortID", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string")).Return(
			// 	models.ShortenEty{
			// 		ShortID:     "3",
			// 		OriginalURL: "https://test.com",
			// 	}, nil).Once()

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
				body:        `{"result": "http://localhost:8080/2"}`,
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
			repo := &mocks.Repositorier{}
			shortURL := services.NewShortURL(cfg.BaseURL, repo)
			h := NewHandler(shortURL, timeout).ApiShortenURL()

			repo.On("Save", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.ShortenEty")).Return(nil)

			h.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

		})
	}
}
