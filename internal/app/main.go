package app

import (
	"fmt"

	"github.com/Mycunycu/shortener/internal/server"
)

func Run() error {
	srv := server.NewServer(":8080")
	err := srv.Run()
	if err != nil {
		return fmt.Errorf("server run err: %v", err)
	}

	return nil
}
