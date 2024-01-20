package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type FileStorage struct {
	inMemoryData *MemoryStorage
	filePath     string
}

func newFileStorage(filePath string) *FileStorage {
	inMemoryData := newMemoryStorage()
	fs := FileStorage{
		inMemoryData: inMemoryData,
		filePath:     filePath,
	}

	err := readFromFile(&fs)
	if err != nil {
		panic(err)
	}

	return &fs
}

func readFromFile(fs *FileStorage) error {
	file, err := os.OpenFile(fs.filePath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	data, err := reader.ReadBytes('\n')

	for {
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		itemURL := URLData{}

		unMarshalErr := json.Unmarshal(data, &itemURL)

		if unMarshalErr != nil {
			return unMarshalErr
		}

		fs.inMemoryData.SaveURL(itemURL)

		data, err = reader.ReadBytes('\n')
	}
	return nil
}

func (fs *FileStorage) Save(data URLData) error {
	file, err := os.OpenFile(fs.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	dataToWrite, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = file.Write(append(dataToWrite, '\n'))

	if err != nil {
		panic(err)
	}

	return nil
}

func (fs *FileStorage) SaveURL(data URLData) (*URLData, error) {
	err := fs.Save(data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (fs *FileStorage) GetShortURL(key string) (*URLData, error) {
	data, err := fs.inMemoryData.GetShortURL(key)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs *FileStorage) Ping() error {
	return nil
}
