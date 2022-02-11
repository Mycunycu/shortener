package app

import (
	"context"
	"errors"
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
	"github.com/Mycunycu/shortener/internal/services"
	"github.com/golang-migrate/migrate/v4"
)

func Run() error {
	cfg := config.New()

	db, err := repository.NewDatabase(cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("error db connection: %v", err)
	}
	defer db.Close()

	err = db.Migrate(cfg.MigrationPath)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	storage, err := repository.NewStorage(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("error creating NewStorage: %v", err)
	}

	shortURL := services.NewShortURL(cfg.BaseURL, db, storage)

	handler := handlers.NewHandler(shortURL, time.Duration(cfg.CtxTimeout)*time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.CtxTimeout)*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}

	fmt.Println("Gracefull stopped")

	return nil
}
