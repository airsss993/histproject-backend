package requests

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// Максимальный размер архива - 50 МБ
const maxArchiveSize = 50 * 1024 * 1024

// CreateRequestReq - структура запроса для создания новой заявки
type CreateRequestReq struct {
	Title            string                `form:"title" binding:"required,max=200"`
	Description      string                `form:"description" binding:"required"`
	EventDate        string                `form:"eventDate" binding:"required,datetime=2006-01-02"`
	EventTypeId      int                   `form:"eventTypeId" binding:"required,gt=0"`
	Email            string                `form:"email" binding:"required,email,max=70"`
	TelegramUsername string                `form:"telegramUsername" binding:"required"`
	Archive          *multipart.FileHeader `form:"archive" binding:"required"`
}

// Структура для сохранения в БД
type RequestData struct {
	Title            string
	Description      string
	Email            string
	TelegramUsername string
	ArchiveURL       string
	SiteURL          string
	ScreenshotURL    string
	EventDate        string
	EventTypeId      int
}

// CreateRequset - метод создания заявки на публикацию сайта, которая пришла от пользователя
func CreateRequest(c *gin.Context) {
	// Инициализируем структуру запроса для того чтобы заполнить её данными
	var req CreateRequestReq

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка парсинга JSON: " + err.Error(),
		})
		return
	}

	// Проверяем правильный формат архива - zip
	if !strings.Contains(req.Archive.Filename, ".zip") {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Архив должен иметь расширение zip",
		})
		return
	}

	// Проверяем допустимый размер архива
	if req.Archive.Size > maxArchiveSize {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Файл должен быть не больше 50 МБ",
		})
		return
	}

	// Загружаем архив в бакет и возвращаем его путь

	// Формируем структуру данных для сохранения в БД
	requestData := RequestData{
		Title:            req.Title,
		Description:      req.Description,
		Email:            req.Email,
		TelegramUsername: req.TelegramUsername,
		EventDate:        req.EventDate,
		EventTypeId:      req.EventTypeId,
	}

	// Сохраняем новую заявку в БД
	if err := createRequestSQL(requestData); err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Ошибка создания заявки: " + err.Error(),
		})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}

// createRequestSQL - функция для получения информации об одном объекте из БД по его ID.
func createRequestSQL(data RequestData) error {
	query := `
		INSERT INTO requests 
			(title,
		    description,
			email,
			telegram_username,
			archive_url,
			site_url,
			screenshot_url,
		    event_date,
		    event_type_id)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.DB.Exec(query, data.Title, data.Description, data.Email, data.TelegramUsername, data.ArchiveURL, data.SiteURL, data.ScreenshotURL, data.EventDate, data.EventTypeId)
	if err != nil {
		return fmt.Errorf("ошибка вставки новой записи о заявке в БД: %w", err)
	}

	return nil
}
