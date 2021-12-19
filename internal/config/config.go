package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerAddress string `yaml:"SERVER_ADDRESS" env-default:"localhost:3000"`
	BaseURL       string `yaml:"BASE_URL" env-default:"http://localhost:3000/"`
}

var cfg Config

func Get() Config {
	return cfg
}

func Init() error {
	var once sync.Once
	var err error

	once.Do(func() {
		if err = cleanenv.ReadConfig("config.yml", &cfg); err != nil {
			err = fmt.Errorf("cleanenv.ReadConfig %w", err)
		}
	})

	return err
}
