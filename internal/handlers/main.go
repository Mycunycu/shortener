package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var mu = &sync.Mutex{}
var urls = make(map[string]string)
var baseShortURL = "http://localhost:8080/"
var id int64 = 0

func ProcessURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path != "/" {
			http.Error(w, "Wrong path", http.StatusBadRequest)
			return
		}

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
	case http.MethodGet:
		id := getID(r)
		if id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		mu.Lock()
		resp := urls[id]
		mu.Unlock()

		fmt.Println(resp)

		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Location", resp)
		w.Header().Write(w)
		return
	default:
		http.Error(w, "Only POST and GET requests are allowed!", http.StatusBadRequest)
		return
	}

}

func getID(r *http.Request) string {
	p := strings.Split(r.URL.Path, "/")

	if len(p) > 1 {
		return p[1]
	}

	return ""
}
