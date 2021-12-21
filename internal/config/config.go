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
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
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

		flag.Func("a", "server address for shorten", func(flagValue string) error {
			cfg.ServerAddress = flagValue
			return nil
		})
		flag.Func("b", "base url for expand", func(flagValue string) error {
			cfg.BaseURL = flagValue
			return nil
		})
		// flag.Func("f", "path to storage file", func(flagValue string) error {
		// 	cfg.FileStoragePath = flagValue
		// 	return nil
		// })

		flag.Parse()

		// if !helpers.IsValidAddress(cfg.ServerAddress) || !helpers.IsValidAddress(cfg.BaseURL) {
		// 	cfg.ServerAddress = DefaultServerAddress
		// 	cfg.BaseURL = DefaultBaseURL
		// }
	})
}
