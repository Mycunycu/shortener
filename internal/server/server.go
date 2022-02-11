package server

import (
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		&http.Server{Addr: addr, Handler: handler},
	}
}
