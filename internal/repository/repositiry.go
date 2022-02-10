package repository

import (
	"context"

	"github.com/Mycunycu/shortener/internal/models"
)

type Repositorier interface {
	Save(context.Context, models.ShortenEty) error
	GetByShortID(context.Context, string) (models.ShortenEty, error)
	GetByUserID(context.Context, string) ([]models.ShortenEty, error)
	PingDB(context.Context) error
}
