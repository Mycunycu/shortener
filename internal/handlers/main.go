package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println(string(b))

	case http.MethodGet:
		id := getID(r)
		if id == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		return
	default:
		http.Error(w, "Only POST and GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func getID(r *http.Request) int {
	p := strings.Split(r.URL.Path, "/")

	if len(p) == 2 {
		id, err := strconv.Atoi(p[1])
		if err != nil {
			return 0
		}

		return id
	}

	return 0
}
