package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type FileStorage struct {
	memoryStorage *MemoryStorage
	filePath      string
}

func (fs *FileStorage) DeleteURL(shortURL string, userID string) error {
	return fs.memoryStorage.DeleteURL(shortURL, userID)
}

func (fs *FileStorage) GetURLsByUserID(userID string) ([]URLData, error) {
	return fs.memoryStorage.GetURLsByUserID(userID)
}

func newFileStorage(filePath string) (*FileStorage, error) {
	memoryStorage, err := newMemoryStorage()
	if err != nil {
		return nil, err
	}

	fs := FileStorage{
		memoryStorage: memoryStorage,
		filePath:      filePath,
	}

	err = readFromFile(fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}

func readFromFile(fs FileStorage) error {
	file, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		var sd URLData
		if err := json.Unmarshal(line, &sd); err != nil {
			return err
		}

		_, err := fs.memoryStorage.SaveURL(sd)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) Save() error {
	file, err := os.OpenFile(fs.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, v := range fs.memoryStorage.storageURL {
		line, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = writer.Write(append(line, '\n'))
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) SaveURL(data URLData) (*URLData, error) {
	fs.memoryStorage.SaveURL(data)
	err := fs.Save()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (fs *FileStorage) GetShortURL(key string) (*URLData, error) {
	data, err := fs.memoryStorage.GetShortURL(key)

	fmt.Println(data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs *FileStorage) Ping() error {
	return nil
}
