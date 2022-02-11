package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Mycunycu/shortener/internal/models"
)

var _ Storager = (*Storage)(nil)

type Storage struct {
	urls map[string][]models.ShortenItem
	mu   *sync.RWMutex
	file *os.File
}

func NewStorage(pathToFile string) (*Storage, error) {
	file, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	stor := &Storage{file: file, mu: &sync.RWMutex{}}
	stor.ReadAllFromFile()

	return stor, nil
}

func (s *Storage) SaveInMemory(userID, shortID, originalURL string) {
	s.mu.Lock()
	_, ok := s.urls[userID]
	if !ok {
		s.urls[userID] = make([]models.ShortenItem, 0)
	}

	s.urls[userID] = append(s.urls[userID], models.ShortenItem{
		ShortID:     shortID,
		OriginalURL: originalURL,
	})
	s.mu.Unlock()
}

func (s *Storage) WriteToFile(userID, shortID, originalURL string) {
	s.file.Write([]byte(fmt.Sprintf("%s|", userID)))
	s.file.Write([]byte(fmt.Sprintf("%s|", shortID)))
	s.file.Write([]byte(fmt.Sprintf("%s\n", originalURL)))
}

func (s *Storage) GetByShortID(id string) (string, error) {
	s.mu.RLock()
	for _, list := range s.urls {
		for _, item := range list {
			if item.ShortID == id {
				return item.OriginalURL, nil
			}
		}
	}
	defer s.mu.RUnlock()

	return "", errors.New("no have data")
}

func (s *Storage) ReadAllFromFile() {
	s.urls = make(map[string][]models.ShortenItem)

	scanner := bufio.NewScanner(s.file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ln := strings.Split(scanner.Text(), "\n")
		splited := strings.Split(ln[0], "|")

		s.SaveInMemory(splited[0], splited[1], splited[2])
	}
}

func (s *Storage) GetAllByUserID(id string) []models.ShortenItem {
	list, ok := s.urls[id]
	if !ok {
		return nil
	}

	return list
}
