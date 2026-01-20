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
		Port string
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
		log.Fatal("Error loading .env file")
	}

	if err := setFromEnv(&cfg); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &cfg, nil
}

func setFromEnv(cfg *Config) error {
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Database.DSN = os.Getenv("PG_DSN")

	if cfg.Server.Port == "" {
		return errors.New("SERVER_PORT environment variable is required")
	}
	if cfg.Database.DSN == "" {
		return errors.New("PG_DSN environment variable is required")
	}

	return nil
}
