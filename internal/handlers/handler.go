package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/google/uuid"
)

const (
	secretKey  = "secret"
	cookieName = "userID"
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

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
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

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
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

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
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

func (h *Handler) UserUrlsById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) PingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.repo.PingDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) setCookie(w http.ResponseWriter, name, value string) {
	encryptedValue := h.encryptCookieValue(value, secretKey)

	http.SetCookie(w, &http.Cookie{
		Name:  name,
		Value: encryptedValue,
	})
}

func (h *Handler) encryptCookieValue(value, key string) string {
	byteValue := []byte(value)
	byteKey := []byte(key)

	mac := hmac.New(sha256.New, byteKey)
	mac.Write(byteValue)
	sig := mac.Sum(nil)
	result := append(byteValue, sig...)

	return hex.EncodeToString(result)
}

func (h *Handler) getUserID(r *http.Request) (string, bool) {
	userID, err := h.getUserIDByCookie(r, cookieName, secretKey)
	if err != nil {
		newID, _ := uuid.NewUUID()
		return newID.String(), true
	}

	return userID, false
}

func (h *Handler) getUserIDByCookie(r *http.Request, cookieName, key string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	byteValue, err := hex.DecodeString(cookie.Value)
	if err != nil {
		return "", err
	}

	byteKey := []byte(key)
	randID, _ := uuid.NewUUID()
	lenID := len([]byte(randID.String()))

	gotUserID := byteValue[:lenID]
	gotSign := byteValue[lenID:]

	mac := hmac.New(sha256.New, byteKey)
	mac.Write(gotUserID)
	sig := mac.Sum(nil)

	if hmac.Equal(sig, gotSign) {
		return string(gotUserID), nil
	}

	return "", errors.New("invalid signature")
}
