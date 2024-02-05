package storage

import (
	"errors"
)

type MemoryStorage struct {
	storageURL map[string]URLData
}

func (s *MemoryStorage) DeleteURL(shortURL string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryStorage) GetURLsByUserID(userID string) ([]URLData, error) {
	var result []URLData

	for _, storageURLItem := range s.storageURL {
		if storageURLItem.UUID == userID {
			result = append(result, storageURLItem)
		}
	}

	return result, nil
}

func newMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storageURL: make(map[string]URLData),
	}
}

func (s *MemoryStorage) Ping() error {
	return nil
}

func (s *MemoryStorage) SaveURL(data URLData) (*URLData, error) {
	s.storageURL[data.ShortURL] = data
	return &data, nil
}

func (s *MemoryStorage) GetShortURL(key string) (*URLData, error) {
	if value, ok := s.storageURL[key]; ok {
		return &value, nil
	}
	return nil, errors.New("URL not found")
}
