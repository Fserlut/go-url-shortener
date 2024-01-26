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
	StorageType     string
	DatabaseDSN     string
}

func InitConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.BaseReturnURL, "b", "http://localhost:8080", "address return url")
	flag.StringVar(&cfg.LogLevel, "l", "info", "logger level")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db.json", "file to save urls")
	flag.StringVar(&cfg.StorageType, "s", "file", "storage to use (memory/file/db)")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database dsn")

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

	if envStorageType := os.Getenv("STORAGE_TYPE"); envStorageType != "" {
		cfg.StorageType = envStorageType
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		cfg.DatabaseDSN = envDatabaseDSN
	}

	if cfg.DatabaseDSN != "" {
		cfg.StorageType = "db"
	}

	return cfg
}
