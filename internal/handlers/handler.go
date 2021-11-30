package handlers

import (
	"io"
	"net/http"

	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

const baseShortURL = "http://localhost:8080/"

type Handler struct {
	repo repository.URLRepository
}

func NewHandler(r repository.URLRepository) *Handler {
	return &Handler{repo: r}
}

func (h *Handler) ShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bOrigURL, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(bOrigURL) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sOrigURL := string(bOrigURL)
		validURL := govalidator.IsURL(sOrigURL)
		if !validURL {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := h.repo.Set(sOrigURL)
		resp := baseShortURL + id

		w.Header().Set("content-type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resp))
	}
}

func (h *Handler) ExpandURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		resp, err := h.repo.GetById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}

		w.Header().Set("Location", resp)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
