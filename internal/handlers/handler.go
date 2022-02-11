package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Mycunycu/shortener/internal/helpers"
	"github.com/Mycunycu/shortener/internal/models"
	"github.com/Mycunycu/shortener/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	secretKey  = "secret"
	cookieName = "userID"
)

type Handler struct {
	shortURL services.ShortURLService
	timeout  time.Duration
}

func NewHandler(shortURL services.ShortURLService, timeout time.Duration) *Handler {
	return &Handler{shortURL: shortURL, timeout: timeout}
}

func (h *Handler) ShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
		}

		originalURL := string(body)

		shortURL, err := h.shortURL.ShortenURL(ctx, userID, originalURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}
}

func (h *Handler) ExpandURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
		}

		originalURL, err := h.shortURL.ExpandURL(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *Handler) ApiShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

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

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
		}

		shortURL, err := h.shortURL.ShortenURL(ctx, userID, req.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responce := models.ShortenResponce{Result: shortURL}
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

func (h *Handler) HistoryByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

		userID, isNewID := h.getUserID(r)
		if isNewID {
			h.setCookie(w, cookieName, userID)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		result, err := h.shortURL.GetHistoryByUserID(ctx, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		jsonResult, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResult)
	}
}

func (h *Handler) PingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

		err := h.shortURL.PingDB(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
