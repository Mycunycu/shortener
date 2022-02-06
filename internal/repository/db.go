package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDB(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
