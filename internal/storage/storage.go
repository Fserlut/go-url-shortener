package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/Fserlut/go-url-shortener/internal/config"
	random "github.com/Fserlut/go-url-shortener/internal/utils"
)

type Storage struct {
	URLStorage map[string]string
	cfg        *config.Config
	File       *os.File
}

type URLJson struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Storage) AddURL(url string) string {
	shortURL := random.GetShortURL()
	s.URLStorage[shortURL] = url

	urlToSave := URLJson{
		UUID:        strconv.Itoa(len(s.URLStorage)),
		OriginalURL: url,
		ShortURL:    shortURL,
	}

	data, err := json.Marshal(urlToSave)

	if err != nil {
		panic(err)
	}

	_, err = s.File.Write(append(data, '\n'))

	if err != nil {
		panic(err)
	}

	return shortURL
}

func initOldURLs(file *os.File) map[string]string {
	URLSFromFile := make(map[string]string)

	if file == nil {
		return URLSFromFile
	}

	reader := bufio.NewReader(file)

	data, err := reader.ReadBytes('\n')

	for {
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		itemURL := URLJson{}

		unMarshalErr := json.Unmarshal(data, &itemURL)

		if unMarshalErr != nil {
			panic(unMarshalErr)
		}

		URLSFromFile[itemURL.ShortURL] = itemURL.OriginalURL

		data, err = reader.ReadBytes('\n')
	}

	return URLSFromFile
}

func InitStorage(cfg *config.Config) *Storage {
	//dir, _ := os.Getwd()

	//file, err := os.OpenFile(dir+cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	URLStorage := initOldURLs(file)

	return &Storage{
		URLStorage: URLStorage,
		cfg:        cfg,
		File:       file,
	}
}
