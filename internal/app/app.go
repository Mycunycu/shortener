package app

import (
	"fmt"

	"github.com/Mycunycu/shortener/internal/routes"
	"github.com/Mycunycu/shortener/internal/server"
)

func Run() error {
	port := ":8080"
	r := routes.NewRouter()

	srv := server.NewServer(port, r)

	err := srv.Run()
	if err != nil {
		return fmt.Errorf("server run err: %v", err)
	}

	return nil
}
