package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Mycunycu/shortener/internal/config"
	"github.com/Mycunycu/shortener/internal/helpers"
	"github.com/Mycunycu/shortener/internal/models"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	cfg  config.Config
	repo repository.URLRepository
}

func NewHandler(cfg config.Config, r repository.URLRepository) *Handler {
	return &Handler{cfg: cfg, repo: r}
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
		resp := h.cfg.BaseURL + "/" + id

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

		resp, err := h.repo.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}

		w.Header().Set("Location", resp)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *Handler) Shorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.ShortenRequest

		err := helpers.DecodeJSONBody(w, r, &req)
		if err != nil {
			var br *helpers.BadRequest
			if errors.As(err, &br) {
				http.Error(w, br.Msg, br.Status)
			} else {
				log.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		validURL := govalidator.IsURL(req.URL)
		if !validURL {
			http.Error(w, "Invalid URL field", http.StatusBadRequest)
			return
		}

		id := h.repo.Set(req.URL)
		result := h.cfg.BaseURL + "/" + id
		responce := models.ShortenResponce{Result: result}

		jsonResp, err := json.Marshal(responce)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	}
}
