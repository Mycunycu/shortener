package services

import (
	"context"
	"errors"

	"github.com/Mycunycu/shortener/internal/models"
	"github.com/Mycunycu/shortener/internal/repository"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

var _ ShortURLService = (*ShortURL)(nil)

type ShortURL struct {
	baseURL string
	db      repository.Repositorier
	//storage repository.Storager
}

func NewShortURL(baseURL string, db repository.Repositorier) *ShortURL {
	return &ShortURL{baseURL: baseURL, db: db}
}

func (s *ShortURL) ShortenURL(ctx context.Context, userID, originalURL string) (string, error) {
	isValid := govalidator.IsURL(originalURL)
	if !isValid {
		return "", errors.New("invalid original URL")
	}

	shortID := uuid.NewString()
	shortURL := s.baseURL + "/" + shortID

	ety := models.ShortenEty{
		UserID:      userID,
		ShortID:     shortID,
		OriginalURL: originalURL,
	}

	err := s.db.Save(ctx, ety)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *ShortURL) ExpandURL(ctx context.Context, id string) (string, error) {
	ety, err := s.db.GetByShortID(ctx, id)
	if err != nil {
		return "", err
	}

	return ety.OriginalURL, err
}

func (s *ShortURL) GetHistoryByUserID(ctx context.Context, id string) ([]models.UserHistoryItem, error) {
	history, err := s.db.GetByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, errors.New("can't find")
	}

	result := make([]models.UserHistoryItem, len(history))
	for i, item := range history {
		result[i] = models.UserHistoryItem{
			ShortURL:    s.baseURL + "/" + item.ShortID,
			OriginalURL: item.OriginalURL,
		}
	}

	return result, nil
}

func (s *ShortURL) PingDB(ctx context.Context) error {
	return s.db.PingDB(context.Background())
}

func (s *ShortURL) ShortenBatch(ctx context.Context, userID string, req models.ShortenBatchRequest) ([]models.BatchItemResponse, error) {
	dataToSave := make([]models.ShortenEty, len(req))
	result := make([]models.BatchItemResponse, len(req))

	for i, item := range req {
		shortID := uuid.NewString()

		dataToSave[i] = models.ShortenEty{
			UserID:      userID,
			ShortID:     shortID,
			OriginalURL: item.OriginalURL,
		}

		result[i] = models.BatchItemResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      s.baseURL + "/" + shortID,
		}
	}

	err := s.db.SaveBatch(ctx, dataToSave)
	if err != nil {
		return nil, err
	}

	return result, nil
}
