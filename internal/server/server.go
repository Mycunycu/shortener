package server

import (
	"net/http"

	"github.com/asaskevich/govalidator"
)

type Server struct {
	*http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	addr := "localhost:8080"
	if govalidator.IsPort(port) {
		addr = port
	}
	
	return &Server{
		&http.Server{Addr: addr, Handler: handler},
	}
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}


