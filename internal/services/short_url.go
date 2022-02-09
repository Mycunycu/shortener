package services

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var _ ShortURLService = (*ShortURL)(nil)

type ShortURL struct {
	baseURL string
	//id   int64
	//urls map[string]string
	//mu      *sync.RWMutex
	//storage *os.File
	db *pgxpool.Pool
}

func NewShortURL(db *pgxpool.Pool, baseURL string) (*ShortURL, error) {
	// file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	// if err != nil {
	// 	return nil, err
	// }

	//shortURL := &ShortURL{storage: file, mu: &sync.RWMutex{}, db: db}
	// storedData := shortURL.ReadAllData()
	// shortURL.id = int64(len(storedData))
	// shortURL.urls = storedData

	shortURL := &ShortURL{db: db, baseURL: baseURL}
	return shortURL, nil
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
