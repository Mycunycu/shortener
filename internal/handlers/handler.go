package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Mycunycu/shortener/internal/helpers"
	"github.com/Mycunycu/shortener/internal/models"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	baseURL string
	repo    repository.Repositorier
}

func NewHandler(baseURL string, r repository.Repositorier) *Handler {
	return &Handler{baseURL: baseURL, repo: r}
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
		h.repo.WriteData(fmt.Sprintf("%s-", id))
		h.repo.WriteData(fmt.Sprintf("%s\n", sOrigURL))
		resp := h.baseURL + "/" + id

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
				return
			}

			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		validURL := govalidator.IsURL(req.URL)
		if !validURL {
			http.Error(w, "Invalid URL field", http.StatusBadRequest)
			return
		}

		id := h.repo.Set(req.URL)
		h.repo.WriteData(fmt.Sprintf("%s-", id))
		h.repo.WriteData(fmt.Sprintf("%s\n", req.URL))
		result := h.baseURL + "/" + id
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
