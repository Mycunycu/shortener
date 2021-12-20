package server

import (
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Server struct {
	*http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	serverAddress := "localhost:8080"
	if isValidAddress(addr) {
		serverAddress = addr
	}

	return &Server{
		&http.Server{Addr: serverAddress, Handler: handler},
	}
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}

func isValidAddress(addr string) bool {
	return govalidator.IsPort(strings.Split(addr, ":")[1])
}
