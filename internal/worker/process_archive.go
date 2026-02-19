package worker

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/pkg/storage"
	"github.com/chromedp/chromedp"
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
		_ = w.updateRequestStatus(payload.RequestId, "Отклонена", "Обнаружен вирус", "", "")
		return nil
	}

	// 4. Распаковка архива и сохранение файлов в хранилище
	if err := w.unzipArchive(payload.RequestId, payload.ArchiveId); err != nil {
		_ = w.updateRequestStatus(payload.RequestId, "Отклонена", "Отсутствует index.html", "", "")
		return nil
	}

	// 5. Создание скриншота сайта и сохранение его в хранилище, а URL в БД
	screenshotUrl, err := w.createScreenshot(payload.RequestId)
	if err != nil {
		return fmt.Errorf("ошибка при создании скриншота: %w", err)
	}

	// 6. Обновление статуса заявки в БД и создание URL сайта
	siteUrl := fmt.Sprintf("http://%s/sites/%d/index.html", w.cfg.Storage.MinioPublicUrl, payload.RequestId)
	_ = w.updateRequestStatus(payload.RequestId, "На модерации", "", siteUrl, screenshotUrl)

	return nil
}
func (w *Worker) createScreenshot(requestID int) (string, error) {
	// Получаем актуальный WebSocket URL от Chrome
	chromeAddr := strings.Replace(w.cfg.App.ChromeUrl, "ws://", "", 1)
	req, err := http.NewRequest("GET", "http://"+chromeAddr+"/json/version", nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса к Chrome: %w", err)
	}
	req.Host = "localhost"

	// Выполняем запрос к Chrome для получения WebSocket URL
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка подключения к Chrome: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Разбираем JSON-ответ от Chrome для получения webSocketDebuggerUrl
	var result map[string]string
	json.Unmarshal(body, &result)

	// Модифицируем webSocketDebuggerUrl, чтобы использовать имя хоста "chrome" и порт 9222
	parsedWsURL, _ := url.Parse(result["webSocketDebuggerUrl"])
	parsedWsURL.Host = "chrome:9222"
	wsURL := parsedWsURL.String()

	// 1. Создание контекста для chromedp
	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		wsURL,
		chromedp.NoModifyURL,
	)
	defer cancel()

	// 2. Создание контекста для выполнения задач chromedp
	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	var buf []byte
	siteUrl := fmt.Sprintf("http://%s/sites/%d/index.html", w.cfg.Storage.MinioEndpoint, requestID)

	// 3. Навигация к сайту и создание скриншота
	if err := chromedp.Run(ctx,
		chromedp.Navigate(siteUrl),
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		return "", fmt.Errorf("ошибка при создании скриншота: %w", err)
	}

	// 4. Загрузка скриншота в хранилище и получение URL
	screenshotKey := fmt.Sprintf("%d/screenshot.png", requestID)
	if err := w.storage.UploadFileSites(bytes.NewReader(buf), int64(len(buf)), screenshotKey); err != nil {
		return "", fmt.Errorf("ошибка при загрузке скриншота в хранилище: %w", err)
	}

	// 5. Формирование URL скриншота
	screenshotUrl := fmt.Sprintf("http://%s/sites/%d/screenshot.png", w.cfg.Storage.MinioPublicUrl, requestID)

	return screenshotUrl, nil
}

// updateRequestStatus - функция для обновления статуса заявки и комментария администратора
func (w *Worker) updateRequestStatus(requestID int, status, comment, siteUrl, screenshorUrl string) error {
	query := `UPDATE requests SET status = $1, admin_comment = $2, site_url = $3, screenshot_url = $4 WHERE id = $5`
	_, err := w.db.Exec(query, status, comment, siteUrl, screenshorUrl, requestID)
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
