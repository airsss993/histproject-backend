package worker

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/pkg/storage"
	"github.com/dutchcoders/go-clamd"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

// Worker - структура для обработки задач в фоновом режиме
type Worker struct {
	db      *sqlx.DB
	storage *storage.MinioClient
	cfg     config.Config
}

// NewWorker - функция для создания нового экземпляра Worker
func NewWorker(db *sqlx.DB, storage *storage.MinioClient, cfg config.Config) *Worker {
	return &Worker{
		db:      db,
		storage: storage,
		cfg:     cfg,
	}
}

// NewMux - функция для создания нового ServeMux и регистрации обработчиков задач
func NewMux(w *Worker) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TypeProcessArchive, w.ProcessArchiveTask)

	return mux
}

// ProcessArchiveTask - обработчик задачи по обработке архива
func (w *Worker) ProcessArchiveTask(ctx context.Context, t *asynq.Task) error {
	// 1. Получение ID заявки и архива
	var payload ProcessArchivePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("ошибка при разборе полезной нагрузки задачи: %w", err)
	}

	log.Printf("[Worker] Начало обработки заявки %d", payload.RequestId)

	// 2. Получение архива из хранилища по ID архива
	archive, _, err := w.storage.GetArchive(payload.ArchiveId)
	if err != nil {
		return fmt.Errorf("ошибка при получении архива из хранилища: %w", err)
	}

	// 3. Проверка архива на вирусы
	if err := w.checkArchiveForViruses(archive); err != nil {
		// Если вирус найден, обновить статус заявки на "ошибка" и добавить комментарий администратора
		_ = w.updateRequestStatus(payload.RequestId, "Отклонена", "Обнаружен вирус", "")
		return nil
	}

	// 4. Распаковка архива и сохранение файлов в хранилище
	if err := w.unzipArchive(payload.RequestId, payload.ArchiveId); err != nil {
		_ = w.updateRequestStatus(payload.RequestId, "Отклонена", "Отсутствует index.html", "")
		return nil
	}

	// TODO: добавить скриншот сайта и сохранить его в хранилище, а URL в БД

	// 5. Обновиление статуса заявки в БД и создание URL сайта
	siteUrl := fmt.Sprintf("http://%s/sites/%d/index.html", w.cfg.Storage.MinioPublicUrl, payload.RequestId)
	_ = w.updateRequestStatus(payload.RequestId, "На модерации", "", siteUrl)

	return nil
}

// updateRequestStatus - функция для обновления статуса заявки и комментария администратора
func (w *Worker) updateRequestStatus(requestID int, status, comment, siteUrl string) error {
	query := `UPDATE requests SET status = $1, admin_comment = $2, site_url = $3 WHERE id = $4`
	_, err := w.db.Exec(query, status, comment, siteUrl, requestID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса заявки: %w", err)
	}
	return nil
}

// unzipArchive - функция для распаковки архива и загрузки файлов в хранилище
func (w *Worker) unzipArchive(requestID int, archiveID string) error {
	// 1. Получить объект из MinIO и его размер
	archive, size, err := w.storage.GetArchive(archiveID)
	if err != nil {
		return fmt.Errorf("ошибка при получении архива из хранилища: %w", err)
	}

	// 2. Открыть как zip через zip.NewReader(object, size)
	zipReader, err := zip.NewReader(archive, size)
	if err != nil {
		return fmt.Errorf("ошибка при открытии архива как zip: %w", err)
	}

	// 3. Проверить наличие index.html в корне или на первом уровне вложенности
	hasIndex := false
	var prefix string
	for _, file := range zipReader.File {
		name := filepath.Base(file.Name)
		parts := strings.Split(strings.Trim(file.Name, "/"), "/")
		if name == "index.html" && len(parts) <= 2 {
			hasIndex = true
			prefix = filepath.Dir(file.Name)
			if prefix == "." {
				prefix = ""
			}
			break
		}
	}
	if !hasIndex {
		return fmt.Errorf("ошибка при проверке наличия index.html в архиве: файл index.html не найден в корне или на первом уровне вложенности")
	}

	// 4. Распаковать все файлы и загрузить их в хранилище, сохраняя структуру папок
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("ошибка открытия файла: %v", err)
		}

		filePath := file.Name
		if prefix != "" {
			filePath = strings.TrimPrefix(file.Name, prefix+"/")
		}
		key := fmt.Sprintf("%d/%s", requestID, filePath)

		if err := w.storage.UploadFileSites(rc, int64(file.UncompressedSize64), key); err != nil {
			_ = rc.Close()
			return fmt.Errorf("ошибка загрузки файла: %v", err)
		}
		_ = rc.Close()
	}

	return nil
}

// checkArchiveForViruses - функция для проверки архива на вирусы с помощью ClamAV
func (w *Worker) checkArchiveForViruses(reader io.Reader) error {
	// Инициализируем клиента ClamAV
	c := clamd.NewClamd(w.cfg.App.ClamAVHost)
	response, err := c.ScanStream(reader, nil)
	if err != nil {
		return fmt.Errorf("ошибка при сканировании архива на вирусы: %w", err)
	}

	// Проверяем результаты сканирования
	for result := range response {
		if result.Status == clamd.RES_FOUND {
			return fmt.Errorf("вирус найден: %s", result.Description)
		}
	}

	return nil
}
