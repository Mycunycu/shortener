package routes

import (
	"github.com/Mycunycu/shortener/internal/handlers"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	*chi.Mux
}

func NewRouter() *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	repo := repository.NewShortURL()
	h := handlers.NewHandler(repo)

	r.Post("/", h.ShortenURL())
	r.Get("/{id}", h.ExpandURL())

	return &Router{r}
}
