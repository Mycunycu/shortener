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

func NewDatabase(ctx context.Context, connStr string) (*Database, error) {

	pool, err := connectDB(ctx, connStr)
	if err != nil {
		return nil, errors.New("db connection error")
	}

	return &Database{pool}, nil
}

func connectDB(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
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

func (d *Database) GetByShortID(ctx context.Context, id string) (models.ShortenEty, error) {
	sql := "SELECT * FROM shortened WHERE short_id = $1"
	row := d.QueryRow(ctx, sql, id)

	var ety models.ShortenEty
	var etyID int
	err := row.Scan(&etyID, &ety.UserID, &ety.ShortID, &ety.OriginalURL)
	return ety, err
}

func (d *Database) GetByOriginalURL(ctx context.Context, url string) (models.ShortenEty, error) {
	sql := "SELECT * FROM shortened WHERE original_url = $1"
	row := d.QueryRow(ctx, sql, url)

	var ety models.ShortenEty
	var etyID int
	err := row.Scan(&etyID, &ety.UserID, &ety.ShortID, &ety.OriginalURL)
	return ety, err
}

func (d *Database) GetByUserID(ctx context.Context, id string) ([]models.ShortenEty, error) {
	sql := "SELECT * FROM shortened WHERE user_id = $1"
	rows, err := d.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}

	var ety models.ShortenEty
	var etyID int
	history := make([]models.ShortenEty, 0)

	for rows.Next() {
		err = rows.Scan(&etyID, &ety.UserID, &ety.ShortID, &ety.OriginalURL)
		if err != nil {
			return nil, err
		}

		history = append(history, ety)
	}

	return history, err
}

func (d *Database) SaveBatch(ctx context.Context, data []models.ShortenEty) error {
	tx, err := d.Begin(ctx)
	if err != nil {
		return err
	}

	sql := "INSERT INTO shortened VALUES (default, $1, $2, $3)"
	stmt, err := tx.Prepare(ctx, "SaveBatch", sql)
	if err != nil {
		return err
	}

	for _, ety := range data {
		_, err := tx.Exec(ctx, stmt.Name, ety.UserID, ety.ShortID, ety.OriginalURL)
		if err != nil {
			return tx.Rollback(ctx)
		}
	}

	return tx.Commit(ctx)
}
