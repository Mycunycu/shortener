package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
)

var mu = &sync.RWMutex{}
var urls = make(map[string]string)
var baseShortURL = "http://localhost:8080/"
var id int64 = 0

func ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	id := atomic.AddInt64(&id, 1)
	idString := strconv.Itoa(int(id))

	mu.Lock()
	urls[idString] = string(bOrigURL)
	mu.Unlock()

	resp := baseShortURL + idString

	w.Header().Set("content-type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))
}

func ExpandURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	mu.RLock()
	resp := urls[id]
	mu.RUnlock()

	fmt.Println(resp)

	w.Header().Set("Location", resp)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
