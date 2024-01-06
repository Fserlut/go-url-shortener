package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseReturnURL   string
	LogLevel        string
	FileStoragePath string
}

func InitConfig() *Config {
	cfg := Config{
		ServerAddress:   ":8080",
		BaseReturnURL:   "http://localhost:8080",
		LogLevel:        "info",
		FileStoragePath: "/tmp/short-url-db.json",
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "address and port to run server")
	flag.StringVar(&cfg.BaseReturnURL, "b", cfg.BaseReturnURL, "address return url")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "logger level")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")

	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		cfg.BaseReturnURL = envServerAddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseReturnURL = envBaseURL
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	return &cfg
}
