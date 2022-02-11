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
	storage repository.Storager
}

func NewShortURL(baseURL string, db repository.Repositorier, stor repository.Storager) *ShortURL {
	return &ShortURL{baseURL: baseURL, db: db, storage: stor}
}

func (s *ShortURL) ShortenURL(userID, originalURL string) (string, error) {
	isValid := govalidator.IsURL(originalURL)
	if !isValid {
		return "", errors.New("invalid original URL")
	}

	shortID := uuid.NewString()
	shortURL := s.baseURL + "/" + shortID

	s.storage.SaveInMemory(userID, shortID, originalURL)
	s.storage.WriteToFile(userID, shortID, originalURL)

	// ety := models.ShortenEty{
	// 	UserID:      userID,
	// 	ShortID:     shortID,
	// 	OriginalURL: originalURL,
	// }

	// err := s.db.Save(ctx, ety)
	// if err != nil {
	// 	return "", err
	// }

	return shortURL, nil
}

func (s *ShortURL) ExpandURL(id string) (string, error) {
	//ety, err := s.db.GetByShortID(ctx, id)
	return s.storage.GetByShortID(id)
}

func (s *ShortURL) GetHistoryByUserID(id string) ([]models.UserHistoryItem, error) {
	history := s.storage.GetAllByUserID(id)
	if history == nil {
		return nil, errors.New("not found")
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
	return nil
	//return s.db.PingDB(context.Background())
}
