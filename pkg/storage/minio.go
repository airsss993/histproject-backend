package storage

import (
	"context"
	"fmt"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minioClient

type minioClient struct {
	cfg config.Storage
	mc  *minio.Client
}

// NewMinioClient создает новый экземпляр Minio Client
func NewMinioClient(cfg config.Storage) *minioClient {
	return &minioClient{
		cfg: cfg,
	}
}

// InitMinio подключается к Minio и создает бакет, если не существует
func (m *minioClient) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(m.cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.cfg.MinioUsername, m.cfg.MinioPassword, ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("ошибка подключения к MinIO: %w", err)
	}

	// Установка подключения Minio
	m.mc = client

	// Проверка наличия бакета и его создание, если не существует
	exists, err := m.mc.BucketExists(ctx, m.cfg.MinioBucketName)
	if err != nil {
		return fmt.Errorf("ошибка проверки наличия бакета: %w", err)
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, m.cfg.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("ошибка создания бакета: %w", err)
		}
	}

	return nil
}
