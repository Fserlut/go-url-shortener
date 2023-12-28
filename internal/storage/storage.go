package storage

import (
	random "github.com/Fserlut/go-url-shortener/internal/utils"
)

type Storage struct {
	URlStorage map[string]string
}

func (s *Storage) AddUrl(url string) string {
	shortURL := random.GetShortURL()
	s.URlStorage[shortURL] = url
	return shortURL
}

func InitStorage() *Storage {
	return &Storage{
		URlStorage: make(map[string]string),
	}
}
