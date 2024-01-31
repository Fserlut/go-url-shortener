package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type ErrURLExists struct{}

func (e *ErrURLExists) Error() string {
	return "This URL already exists"
}

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
        user_id TEXT NOT NULL,
        short_url TEXT NOT NULL UNIQUE,
        original_url TEXT NOT NULL
    );

		CREATE UNIQUE INDEX  IF NOT EXISTS links_original_url_uniq_index
		    on links (original_url);
	`)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &DatabaseStorage{
		db: db,
	}
}

func (s *DatabaseStorage) SaveURL(data URLData) (*URLData, error) {
	res, err := s.db.ExecContext(
		context.Background(),
		`INSERT INTO links (uuid, user_id, short_url, original_url) VALUES ($1, $2, $3, $4) ON CONFLICT (original_url) DO NOTHING`, data.UUID, data.ShortURL, data.OriginalURL, data.UserID,
	)

	if err != nil {
		return nil, err
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		var (
			uuid        string
			shortURL    string
			originalURL string
		)
		row := s.db.QueryRowContext(
			context.Background(),
			"SELECT uuid, short_url, original_url FROM links WHERE original_url = $1", data.OriginalURL,
		)
		err := row.Scan(&uuid, &shortURL, &originalURL)

		if err != nil {
			return nil, err
		}
		return &URLData{UUID: uuid, ShortURL: shortURL, OriginalURL: originalURL}, &ErrURLExists{}
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

func (s *DatabaseStorage) GetURLsByUserID(userID string) ([]URLData, error) {
	var (
		entity URLData
		result []URLData
	)

	query := "select short_url, original_url from links where user_id=$1"
	rows, err := s.db.Query(query, userID)

	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		err = rows.Scan(&entity.ShortURL, &entity.OriginalURL)
		if err != nil {
			break
		}
		result = append(result, entity)
	}

	return result, nil
}

func (s *DatabaseStorage) Ping() error {
	if err := s.db.PingContext(context.Background()); err != nil {
		return err
	}
	return nil
}
