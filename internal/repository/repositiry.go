package repository

import "context"

type Repositorier interface {
	// Set(url string) string
	// GetByID(id string) (string, error)
	// WriteData(data string) error
	// ReadAllData() map[string]string
	PingDB(context.Context) error
}
