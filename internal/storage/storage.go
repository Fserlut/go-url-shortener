package storage

import (
	random "github.com/Fserlut/go-url-shortener/internal/utils"
)

type Storage struct {
	URLStorage map[string]string
}

func (s *Storage) AddURL(url string) string {
	shortURL := random.GetShortURL()
	s.URLStorage[shortURL] = url
	return shortURL
}

func InitStorage() *Storage {
	return &Storage{
		URLStorage: make(map[string]string),
	}
}
