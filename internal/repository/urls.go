package repository

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type ShortURL struct {
	id      int64
	urls    map[string]string
	mu      *sync.RWMutex
	storage *os.File
}

func NewShortURL(path string) (*ShortURL, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	shortURL := &ShortURL{storage: file, mu: &sync.RWMutex{}}
	storedData := shortURL.ReadAllData()
	shortURL.id = int64(len(storedData))
	shortURL.urls = storedData

	return shortURL, nil
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

func (s *ShortURL) WriteData(data string) error {
	_, err := s.storage.Write([]byte(data))
	return err
}

func (s *ShortURL) ReadAllData() map[string]string {
	result := make(map[string]string)

	scanner := bufio.NewScanner(s.storage)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ln := strings.Split(scanner.Text(), "\n")
		splited := strings.Split(ln[0], "-")

		result[splited[0]] = splited[1]
	}

	return result
}
