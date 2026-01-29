package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      App
		Database Database
		CORS     CORS
	}
	App struct {
		Port        string
		SwaggerHost string
	}

	Database struct {
		DSN string
	}

	CORS struct {
		AllowedOrigins string
	}
)

// Init загружает переменные окружения из .env файла и инициализирует конфигурацию приложения.
func Init() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	if err := setFromEnv(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка получения env переменных: %w", err)
	}

	return &cfg, nil
}

// setFromEnv заполняет структуру конфигурации значениями из переменных окружения и валидирует их наличие.
func setFromEnv(cfg *Config) error {
	cfg.App.Port = os.Getenv("SERVER_PORT")
	cfg.App.SwaggerHost = os.Getenv("SWAGGER_HOST")
	cfg.Database.DSN = os.Getenv("PG_DSN")
	cfg.CORS.AllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")

	if cfg.App.Port == "" {
		return errors.New("SERVER_PORT должно быть указано")
	}
	if cfg.App.SwaggerHost == "" {
		return errors.New("SWAGGER_HOST должно быть указано")
	}
	if cfg.Database.DSN == "" {
		return errors.New("PG_DSN должно быть указано")
	}

	return nil
}
