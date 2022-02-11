package services

import (
	"context"

	"github.com/Mycunycu/shortener/internal/models"
)

type ShortURLService interface {
	ShortenURL(string, string) (string, error)
	ExpandURL(string) (string, error)
	GetHistoryByUserID(string) ([]models.UserHistoryItem, error)
	PingDB(context.Context) error
}
