package server

import (
	"net/http"

	"github.com/Mycunycu/shortener/internal/handlers"
)

type Server struct {
	*http.Server
}

func NewServer(port string) *Server {
	srv := &Server{
		&http.Server{Addr: port},
	}

	return srv
}

func (s *Server) Run() error {
	http.HandleFunc("/", handlers.ShortenURL)
	return s.ListenAndServe()
}
