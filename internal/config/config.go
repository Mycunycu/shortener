package config

import (
	"os"
	"sync"
)

const (
	ServerAddress = "SERVER_ADDRESS"
	BaseURL = "BASE_URL"
)

type Config struct {
	ServerAddress string
	BaseURL       string 
}

var cfg Config

func Get() Config {
	initialize()
	return cfg
}

func initialize() {
	var once sync.Once

	once.Do(func() {
		cfg = Config{
			ServerAddress: getEnv(ServerAddress, "localhost:8080"),
			BaseURL: getEnv(BaseURL, "http://localhost:8080/"),
		}
		// if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		// 	log.Fatalf("read config err: %v", err)
		// }
	})
}

func getEnv(name, def string) string {
	value, ok := os.LookupEnv(name)
	if ok {
		return value
	}

	return def
}
