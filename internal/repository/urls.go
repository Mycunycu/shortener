package repository

import (
	"errors"
	"sync"

	"github.com/Mycunycu/shortener/internal/helpers"
)

type ShortURL struct {
	urls map[string]string
	mu   *sync.RWMutex
}

func NewShortURL() *ShortURL {
	return &ShortURL{
		urls: make(map[string]string),
		mu:   &sync.RWMutex{},
	}
}

func (s *ShortURL) Set(url string) string {
	id := helpers.CreateNewId(5)

	s.mu.Lock()
	s.urls[id] = url
	s.mu.Unlock()

	return id
}
func (s *ShortURL) GetByID(id string) (string, error) {
	s.mu.RLock()
	url, ok := s.urls[id]
	s.mu.RUnlock()
	if !ok {
		return "", errors.New("no have data")
	}

	return url, nil
}
