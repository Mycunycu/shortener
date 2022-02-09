package repository

import (
	"context"

	"github.com/Mycunycu/shortener/internal/models"
)

type Repositorier interface {
	// Set(url string) string
	// GetByID(id string) (string, error)
	// WriteData(data string) error
	// ReadAllData() map[string]string
	Save(context.Context, models.ShortenEty) error
	GetByShortID(context.Context, string) (*models.ShortenEty, error)
	PingDB(context.Context) error
}
