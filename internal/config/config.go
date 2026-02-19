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
		Database Database
		CORS     CORS
		Storage  Storage
		Redis    Redis
	}
	App struct {
		Port        string
		SwaggerHost string
		ClamAVHost  string
		ChromeUrl   string
	}

	Database struct {
		DSN string
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
	cfg.App.ClamAVHost = os.Getenv("CLAMAV_HOST")
	cfg.App.ChromeUrl = os.Getenv("CHROME_URL")

	cfg.Database.DSN = os.Getenv("PG_DSN")
	cfg.CORS.AllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")

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
	if cfg.App.SwaggerHost == "" {
		return errors.New("SWAGGER_HOST должно быть указано")
	}
	if cfg.App.ClamAVHost == "" {
		return errors.New("CLAMAV_HOST должно быть указано")
	}

	if cfg.Database.DSN == "" {
		return errors.New("PG_DSN должно быть указано")
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
