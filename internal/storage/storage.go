package storage

import "github.com/Fserlut/go-url-shortener/internal/config"

type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	SaveURL(data URLData) (*URLData, error)
	GetShortURL(key string) (*URLData, error)
	Ping() error
}

func NewStorage(cfg *config.Config) Storage {
	switch cfg.StorageType {
	case "memory":
		return newMemoryStorage()
	case "file":
		return newFileStorage(cfg.FileStoragePath)
	case "db":
		return newDBStorage(cfg.DatabaseDSN)
	default:
		return newMemoryStorage()
	}
}
