package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      App
		Auth     Auth
		Database Database
		CORS     CORS
		Storage  Storage
		Redis    Redis
		N8N      N8N
	Telegram Telegram
	}
	App struct {
		Port          string
		ClamAVEnabled bool
		ClamAVHost    string
		ChromeUrl     string
	}
	Auth struct {
		JWTSecret string
	}

	Database struct {
		DSN string
		Env string
	}

	Storage struct {
		MinioUsername   string
		MinioPassword   string
		MinioEndpoint   string
		MinioBucketName string
		MinioPublicUrl  string
	}

	Redis struct {
		RedisPassword string
		RedisAddress  string
		RedisDb       int
	}

	CORS struct {
		AllowedOrigins string
	}

	N8N struct {
		WebhookURL    string
		WebhookSecret string
	}

	Telegram struct {
		BotToken   string
		WebhookURL string
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
	cfg.App.ClamAVEnabled = os.Getenv("CLAMAV_ENABLED") == "true"
	cfg.App.ClamAVHost = os.Getenv("CLAMAV_HOST")
	cfg.App.ChromeUrl = os.Getenv("CHROME_URL")

	cfg.Auth.JWTSecret = os.Getenv("JWT_SECRET")

	cfg.Database.Env = os.Getenv("APP_ENV")
	switch cfg.Database.Env {
	case "prod":
		cfg.Database.DSN = os.Getenv("PG_DSN_PROD")
	case "dev":
		cfg.Database.DSN = os.Getenv("PG_DSN_DEV")
	default:
		cfg.Database.DSN = os.Getenv("PG_DSN")
	}
	cfg.CORS.AllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")
	cfg.N8N.WebhookURL = os.Getenv("N8N_WEBHOOK_URL")
	cfg.N8N.WebhookSecret = os.Getenv("N8N_WEBHOOK_SECRET")

	cfg.Telegram.BotToken = os.Getenv("TG_BOT_TOKEN")
	cfg.Telegram.WebhookURL = os.Getenv("TG_WEBHOOK_URL")

	cfg.Storage.MinioEndpoint = os.Getenv("MINIO_ENDPOINT")
	cfg.Storage.MinioUsername = os.Getenv("MINIO_ROOT_USER")
	cfg.Storage.MinioPassword = os.Getenv("MINIO_ROOT_PASSWORD")
	cfg.Storage.MinioBucketName = os.Getenv("MINIO_BUCKET_NAME")
	cfg.Storage.MinioPublicUrl = os.Getenv("MINIO_PUBLIC_URL")

	cfg.Redis.RedisPassword = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.RedisAddress = os.Getenv("REDIS_ADDRESS")
	cfg.Redis.RedisDb, _ = strconv.Atoi(os.Getenv("REDIS_DB"))

	if cfg.App.Port == "" {
		return errors.New("SERVER_PORT должно быть указано")
	}
	if cfg.Auth.JWTSecret == "" {
		return errors.New("JWT_SECRET должно быть указано")
	}
	if cfg.App.ClamAVEnabled && cfg.App.ClamAVHost == "" {
		return errors.New("CLAMAV_HOST должно быть указано если CLAMAV_ENABLED=true")
	}

	if cfg.Database.DSN == "" {
		switch cfg.Database.Env {
		case "prod":
			return errors.New("PG_DSN_PROD должно быть указано при APP_ENV=prod")
		case "dev":
			return errors.New("PG_DSN_DEV должно быть указано при APP_ENV=dev")
		default:
			return errors.New("PG_DSN должно быть указано (или укажите APP_ENV=dev/prod и соответствующий DSN)")
		}
	}
	if cfg.Storage.MinioEndpoint == "" {
		return errors.New("MINIO_ENDPOINT должно быть указано")
	}
	if cfg.Storage.MinioUsername == "" {
		return errors.New("MINIO_ROOT_USER должно быть указано")
	}
	if cfg.Storage.MinioPassword == "" {
		return errors.New("MINIO_ROOT_PASSWORD должно быть указано")
	}
	if cfg.Storage.MinioBucketName == "" {
		return errors.New("MINIO_BUCKET_NAME должно быть указано")
	}

	if cfg.Redis.RedisPassword == "" {
		return errors.New("REDIS_PASSWORD должно быть указано")
	}
	if cfg.Redis.RedisAddress == "" {
		return errors.New("REDIS_ADDRESS должно быть указано")
	}

	return nil
}
