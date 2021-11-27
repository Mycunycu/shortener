package server

import (
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	srv := &Server{
		&http.Server{Addr: port, Handler: handler},
	}

	return srv
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}
