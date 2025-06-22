package postgres

import (
	"context"
	"fmt"

	"github.com/arvinloc/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS url(
        id SERIAL PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = db.Exec(context.Background(), `CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.postgres.SaveURL"
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO url(alias,url) VALUES($1, $2)", alias, urlToSave)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var result string

	err := s.db.QueryRow(context.Background(), "SELECT url FROM url WHERE alias = $1", alias).Scan(&result)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"
	result, err := s.db.Exec(context.Background(), "DELETE FROM url WHERE alias = $1", alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return storage.ErrURLNotFound
	}

	return nil
}
