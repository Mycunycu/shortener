package repository

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

type ShortURL struct {
	id   int64
	urls map[string]string
	mu   *sync.RWMutex
}

func NewShortURL() *ShortURL {
	return &ShortURL{
		id:   0,
		urls: make(map[string]string),
		mu:   &sync.RWMutex{},
	}
}

func (s *ShortURL) Set(url string) string {
	atomic.AddInt64(&s.id, 1)
	idString := strconv.Itoa(int(s.id))

	s.mu.Lock()
	s.urls[idString] = url
	s.mu.Unlock()

	return idString
}
func (s *ShortURL) GetById(id string) (string, error) {
	s.mu.RLock()
	url, ok := s.urls[id]
	s.mu.RUnlock()
	if !ok {
		return "", errors.New("no have data")
	}

	return url, nil
}
