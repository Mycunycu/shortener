package config

import (
	"flag"
	"log"
	"sync"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8000"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8000"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"internal/repository/storage.txt"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"postgres://user:123@localhost:5432/practicum?sslmode=disable"`
	MigrationPath   string `env:"MIGRATION_PATH" envDefault:"file://internal/repository/migrations"`
}

var cfg Config

func New() Config {
	initialize()
	return cfg
}

func initialize() {
	var once sync.Once

	once.Do(func() {
		cfg = Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("initialize config error: %v", err)
		}

		flag.Func("a", "server address", func(value string) error {
			cfg.ServerAddress = value
			return nil
		})
		flag.Func("b", "base url", func(value string) error {
			cfg.BaseURL = value
			return nil
		})
		flag.Func("f", "path to storage file", func(value string) error {
			cfg.FileStoragePath = value
			return nil
		})
		flag.Func("d", "database url", func(value string) error {
			cfg.FileStoragePath = value
			return nil
		})

		flag.Parse()
	})
}
