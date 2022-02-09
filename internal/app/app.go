package app

import (
	"fmt"
	"log"
	"os"

	"github.com/Mycunycu/shortener/internal/config"
	"github.com/Mycunycu/shortener/internal/handlers"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/Mycunycu/shortener/internal/routes"
	"github.com/Mycunycu/shortener/internal/server"
	"github.com/Mycunycu/shortener/internal/services"
)

func Run() error {
	cfg := config.New()

	db, err := repository.ConnectDB(cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("error db connection: %v", err)
	}

	shortURL, err := services.NewShortURL(db, cfg.BaseURL)
	if err != nil {
		return fmt.Errorf("error creating new NewShortURL: %v", err)
	}

	handler := handlers.NewHandler(shortURL)
	router := routes.NewRouter(handler)
	srv := server.NewServer(cfg.ServerAddress, router)

	go func() {
		err := srv.Run()
		if err != nil {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer func() {
	// 	cancel()
	// }()

	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatalf("server shutdown failed:%+v", err)
	// }

	// fmt.Println("Gracefull stopped")

	return nil
}
