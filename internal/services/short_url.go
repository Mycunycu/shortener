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
	db      *repository.Database
}

func NewShortURL(baseURL string, db *repository.Database) (*ShortURL, error) {
	shortURL := &ShortURL{baseURL: baseURL, db: db}
	return shortURL, nil
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
	return ety.OriginalURL, err
}

func (s *ShortURL) GetHistoryByUserID(ctx context.Context, id string) ([]models.UserHistoryItem, error) {
	history, err := s.db.GetByUserID(ctx, id)
	if err != nil {
		return nil, err
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

// func (s *ShortURL) Set(url string) string {
// 	atomic.AddInt64(&s.id, 1)
// 	idString := strconv.Itoa(int(s.id))

// 	s.mu.Lock()
// 	s.urls[idString] = url
// 	s.mu.Unlock()

// 	return idString
// }

// func (s *ShortURL) GetByID(id string) (string, error) {
// 	s.mu.RLock()
// 	url, ok := s.urls[id]
// 	s.mu.RUnlock()
// 	if !ok {
// 		return "", errors.New("no have data")
// 	}

// 	return url, nil
// }

// func (s *ShortURL) WriteData(data string) error {
// 	_, err := s.storage.Write([]byte(data))
// 	return err
// }

// func (s *ShortURL) ReadAllData() map[string]string {
// 	result := make(map[string]string)

// 	scanner := bufio.NewScanner(s.storage)
// 	scanner.Split(bufio.ScanLines)

// 	for scanner.Scan() {
// 		ln := strings.Split(scanner.Text(), "\n")
// 		splited := strings.Split(ln[0], "-")

// 		result[splited[0]] = splited[1]
// 	}

// 	return result
// }

func (s *ShortURL) PingDB(ctx context.Context) error {
	return s.db.Ping(context.Background())
}
