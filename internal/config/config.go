package config

import (
	"flag"
	"log"
	"sync"

	"github.com/caarlos0/env"
)

const (
	DefaultServerAddress = ":8080"
	DefaultBaseURL       = "http://localhost:8080/"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.txt"`
}

var cfg Config

func Get() Config {
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

		flag.Parse()
	})
}
