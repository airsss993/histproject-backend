package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// UploadArchive - функция загрузки архива в бакет MinIO
func (m *minioClient) UploadArchive(file *multipart.FileHeader) (string, error) {
	objectId := uuid.New().String()

	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}
	defer fileReader.Close()

	// Загрузка файла в бакет
	_, err = m.mc.PutObject(context.Background(), m.cfg.MinioBucketName, objectId, fileReader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("ошибка при загрузки файла в бакет %s: %v", file.Filename, err)
	}

	// Возвращаем ID объекта
	return objectId, nil
}
