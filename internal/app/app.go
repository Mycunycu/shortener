package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mycunycu/shortener/internal/config"
	"github.com/Mycunycu/shortener/internal/handlers"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/Mycunycu/shortener/internal/routes"
	"github.com/Mycunycu/shortener/internal/server"
)

func Run() error {
	cfg := config.Get()

	storage, err := repository.NewStorage(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("error creating new storage: %v", err)
	}

	repo := repository.NewShortURL(storage)
	handler := handlers.NewHandler(cfg, repo, storage)
	router := routes.NewRouter(handler)
	srv := server.NewServer(cfg.ServerAddress, router)

	go func() {
		err := srv.Run()
		if err != nil {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}

	fmt.Println("Gracefull stopped")

	return nil
}
