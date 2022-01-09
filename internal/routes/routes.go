package routes

import (
	"github.com/Mycunycu/shortener/internal/handlers"
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
	r.Use(middleware.Compress(5))

	r.Post("/", h.ShortenURL())
	r.Get("/{id}", h.ExpandURL())

	r.Post("/api/shorten", h.Shorten())

	return &Router{r}
}
