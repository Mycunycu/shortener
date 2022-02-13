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
	SaveBatch(context.Context, []models.ShortenEty) error
}

type Storager interface {
	SaveInMemory(string, string, string)
	WriteToFile(string, string, string)
	GetByShortID(string) (string, error)
	ReadAllFromFile()
	GetAllByUserID(string) []models.ShortenItem
}
