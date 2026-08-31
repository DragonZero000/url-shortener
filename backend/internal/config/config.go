package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	WebBaseURL  string
}

func Load() (*Config, error) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		slog.Info("no .env file found, using system env")
	}
	var cfg Config
	databaseURL, exist := os.LookupEnv("DATABASE_URL")
	if !exist {
		err = errors.New("no database url  in config")
		slog.Error("no database url", "error", err)
		return nil, err
	}
	cfg.DatabaseURL = databaseURL
	serverPort, exist := os.LookupEnv("SERVER_PORT")
	if !exist {
		serverPort = ":8080"
	}
	cfg.ServerPort = serverPort
	webBaseURL, exist := os.LookupEnv("BASE_URL")
	if !exist {
		err = errors.New("no domain name in config")
		slog.Error("no domain name", "error", err)
		return nil, err
	}
	cfg.WebBaseURL = webBaseURL
	return &cfg, nil
}
