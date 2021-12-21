package repository

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

type ShortURL struct {
	id      int64
	urls    map[string]string
	mu      *sync.RWMutex
	storage IStorage
}

func NewShortURL(storage IStorage) *ShortURL {
	storedData := storage.ReadAll()

	fmt.Println(storedData)

	return &ShortURL{
		id:      int64(len(storedData)),
		urls:    storedData,
		mu:      &sync.RWMutex{},
		storage: storage,
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

func (s *ShortURL) GetByID(id string) (string, error) {
	s.mu.RLock()
	url, ok := s.urls[id]
	s.mu.RUnlock()
	if !ok {
		return "", errors.New("no have data")
	}

	return url, nil
}
