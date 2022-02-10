package services

import (
	"context"

	"github.com/Mycunycu/shortener/internal/models"
)

type ShortURLService interface {
	ShortenURL(context.Context, string, string) (string, error)
	ExpandURL(context.Context, string) (string, error)
	GetHistoryByUserID(context.Context, string) ([]models.UserHistoryItem, error)
	PingDB(context.Context) error
}