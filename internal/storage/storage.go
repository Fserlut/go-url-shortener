package storage

import (
	"fmt"
	"math/rand"
	"time"
)

type Storage struct {
	URlStorage map[string]string
}

func getShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2:8]
}

func (s *Storage) AddUrl(url string) {
	shortURL := getShortURL()
	s.URlStorage[shortURL] = url
}

func InitStorage() *Storage {
	return &Storage{
		URlStorage: make(map[string]string),
	}
}
