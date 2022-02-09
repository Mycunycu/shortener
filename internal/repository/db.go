package repository

import (
	"context"
	"errors"

	"github.com/Mycunycu/shortener/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ Repositorier = (*Database)(nil)

type Database struct {
	*pgxpool.Pool
}

func NewDatabase(connStr string) (*Database, error) {
	pool, err := connectDB(connStr)
	if err != nil {
		return nil, errors.New("db connection error")
	}

	return &Database{pool}, nil
}

func connectDB(connStr string) (*pgxpool.Pool, error) {
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

func (d *Database) Migrate(source string) error {
	m, err := migrate.New(source, d.Config().ConnString())
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func (d *Database) PingDB(ctx context.Context) error {
	return d.Ping(ctx)
}

func (d *Database) Save(ctx context.Context, e models.ShortenEty) error {
	sql := "INSERT INTO shortened VALUES (default, $1, $2, $3)"
	_, err := d.Exec(ctx, sql, e.UserID, e.ShortID, e.OriginalURL)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetByShortID(ctx context.Context, id string) (*models.ShortenEty, error) {
	sql := "SELECT * FROM shortened WHERE short_id = $1"
	row := d.QueryRow(ctx, sql, id)

	var ety models.ShortenEty
	var etyID int
	err := row.Scan(&etyID, &ety.UserID, &ety.ShortID, &ety.OriginalURL)

	return &ety, err
}
