package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   Server
		Database Database
		App      App
	}
	Server struct {
		Port     string
		BasePath string
	}

	Database struct {
		DSN string
	}

	App struct {
		WebURL string
	}
)

func Init() (*Config, error) {
	var cfg Config

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	if err := setFromEnv(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка получения env переменных: %w", err)
	}

	return &cfg, nil
}

func setFromEnv(cfg *Config) error {
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Server.BasePath = os.Getenv("BASE_PATH")
	cfg.Database.DSN = os.Getenv("PG_DSN")

	if cfg.Server.Port == "" {
		return errors.New("SERVER_PORT должно быть указано")
	}
	if cfg.Server.BasePath == "" {
		return errors.New("BASE_PATH должно быть указано")
	}
	if cfg.Database.DSN == "" {
		return errors.New("PG_DSN должно быть указано")
	}

	return nil
}
