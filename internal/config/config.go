package config

import (
	"os"
	"sync"

	"github.com/asaskevich/govalidator"
)

const (
	ServerAddress = "SERVER_ADDRESS"
	BaseURL       = "BASE_URL"
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
		var address = "localhost:8080"
		if isValid := govalidator.IsPort(os.Getenv(ServerAddress)); isValid {
			address = os.Getenv(ServerAddress)
		}

		var baseURL = "http://localhost:8080/"
		if isValid := govalidator.IsURL(os.Getenv(BaseURL)); isValid {
			baseURL = os.Getenv(BaseURL)
		}

		cfg = Config{
			ServerAddress: address,
			BaseURL:       baseURL,
		}
	})
}

// func getEnv(name, def string) string {
// 	value, ok := os.LookupEnv(name)
// 	if ok {
// 		return value
// 	}

// 	return def
// }
