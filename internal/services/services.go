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
	ShortenBatch(context.Context, string, models.ShortenBatchRequest) ([]models.BatchItemResponse, error)
	DeleteBatch(context.Context, string, []string) error
}
