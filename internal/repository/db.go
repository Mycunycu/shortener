package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func ConnectDB(connStr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
