package server

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(ctx context.Context, addr string, handler http.Handler) *Server {
	return &Server{
		&http.Server{Addr: addr, Handler: handler, BaseContext: func(_ net.Listener) context.Context { return ctx }},
	}
}
