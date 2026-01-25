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
	}
	Server struct {
		Port string
	}

	Database struct {
		DSN string
	}
)

// Init загружает переменные окружения из .env файла и инициализирует конфигурацию приложения.
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

// setFromEnv заполняет структуру конфигурации значениями из переменных окружения и валидирует их наличие.
func setFromEnv(cfg *Config) error {
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Database.DSN = os.Getenv("PG_DSN")

	if cfg.Server.Port == "" {
		return errors.New("SERVER_PORT должно быть указано")
	}
	if cfg.Database.DSN == "" {
		return errors.New("PG_DSN должно быть указано")
	}

	return nil
}
