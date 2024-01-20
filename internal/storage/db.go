package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseStorage struct {
	DB *sql.DB
}

func newDBStorage(dsn string) *DatabaseStorage {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &DatabaseStorage{
		DB: db,
	}
}

func (s *DatabaseStorage) SaveURL(data URLData) (*URLData, error) {
	fmt.Println(data)
	return &URLData{}, nil
}

func (s *DatabaseStorage) GetShortURL(key string) (*URLData, error) {
	fmt.Println(key)
	return &URLData{}, nil
}

func (s *DatabaseStorage) Ping() error {
	if err := s.DB.PingContext(context.Background()); err != nil {
		return err
	}
	return nil
}
