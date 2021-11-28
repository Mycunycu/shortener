package server

import (
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		&http.Server{Addr: port, Handler: handler},
	}
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}
