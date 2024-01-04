package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Fserlut/go-url-shortener/internal/config"
	random "github.com/Fserlut/go-url-shortener/internal/utils"
)

type Storage struct {
	URLStorage map[string]string
	cfg        *config.Config
	file       *os.File
}

type URLJson struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func (s *Storage) AddURL(url string) string {
	shortURL := random.GetShortURL()
	s.URLStorage[shortURL] = url
	urlToSave := URLJson{
		Uuid:        strconv.Itoa(len(s.URLStorage)),
		OriginalUrl: url,
		ShortUrl:    shortURL,
	}
	fmt.Println(urlToSave)
	data, err := json.MarshalIndent(urlToSave, "", "\n")

	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	fileWriteErr := os.WriteFile(s.cfg.FileStoragePath, data, 0666)

	if fileWriteErr != nil {
		panic(fileWriteErr)
	}

	return shortURL
}

func InitStorage(cfg *config.Config) *Storage {
	fmt.Println(cfg.FileStoragePath)

	return &Storage{
		URLStorage: make(map[string]string),
		cfg:        cfg,
	}
}
