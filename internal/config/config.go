package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	BaseReturnURL string
}

func InitConfig() *Config {
	cfg := Config{
		ServerAddress: ":8080",
		BaseReturnURL: "http://localhost:8080",
	}
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "address and port to run server")
	flag.StringVar(&cfg.BaseReturnURL, "b", cfg.BaseReturnURL, "address return url")

	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		cfg.BaseReturnURL = envServerAddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseReturnURL = envBaseURL
	}

	return &cfg
}
