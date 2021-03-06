package routes

import (
	"github.com/Mycunycu/shortener/internal/handlers"
	customMiddleware "github.com/Mycunycu/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	*chi.Mux
}

func NewRouter(h *handlers.Handler) *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.GzipCompress)
	r.Use(customMiddleware.GzipDecompress)

	r.Post("/", h.ShortenURL())
	r.Get("/{id}", h.ExpandURL())

	r.Post("/api/shorten", h.APIShortenURL())
	r.Get("/api/user/urls", h.HistoryByUserID())
	r.Delete("/api/user/urls", h.DeleteShortened())
	r.Post("/api/shorten/batch", h.ShortenBatchURL())

	r.Get("/ping", h.PingDB())

	return &Router{r}
}
