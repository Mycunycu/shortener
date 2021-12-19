package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerAddress string `yaml:"SERVER_ADDRESS" env-default:"localhost:3000"`
	BaseURL       string `yaml:"BASE_URL" env-default:"http://localhost:3000/"`
}

var cfg Config

func Get() Config {
	initialize()
	return cfg
}

func initialize() {
	var once sync.Once

	once.Do(func() {
		if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
			log.Fatalf("read config err: %v", err)
		}
	})
}
