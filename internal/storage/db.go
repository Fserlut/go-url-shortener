package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type DatabaseStorage struct {
	db *sql.DB
}

func newDBStorage(dsn string) *DatabaseStorage {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
        uuid TEXT PRIMARY KEY,
        short_url TEXT NOT NULL UNIQUE,
        original_url TEXT NOT NULL
    );
	`)

	if err != nil {
		panic(err)
	}

	return &DatabaseStorage{
		db: db,
	}
}

func (s *DatabaseStorage) SaveURL(data URLData) (*URLData, error) {
	_, err := s.db.ExecContext(
		context.Background(),
		`INSERT INTO links (uuid, short_url, original_url) VALUES ($1, $2, $3)`, data.UUID, data.ShortURL, data.OriginalURL,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &data, nil
}

func (s *DatabaseStorage) GetShortURL(key string) (*URLData, error) {
	var (
		uuid        string
		shortURL    string
		originalURL string
	)
	row := s.db.QueryRowContext(
		context.Background(),
		"SELECT uuid, short_url, original_url FROM links WHERE short_url = $1", key,
	)
	err := row.Scan(&uuid, &shortURL, &originalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("URL not found")
		} else {
			return nil, err
		}
	}

	return &URLData{
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}, nil
}

func (s *DatabaseStorage) Ping() error {
	if err := s.db.PingContext(context.Background()); err != nil {
		return err
	}
	return nil
}
