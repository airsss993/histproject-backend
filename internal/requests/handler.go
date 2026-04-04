package requests

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler — HTTP-обработчики заявок.
type Handler struct {
	svc *Service
}

// NewHandler создаёт handler заявок.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// CreateRequestReq — форма создания заявки.
type CreateRequestReq struct {
	Title            string                `form:"title" binding:"required,max=200"`
	Description      string                `form:"description" binding:"required"`
	EventDate        string                `form:"eventDate" binding:"required,datetime=2006-01-02"`
	EventTypeId      int                   `form:"eventTypeId" binding:"required,gt=0"`
	Email            string                `form:"email" binding:"required,email,max=70"`
	TelegramUsername string                `form:"telegramUsername" binding:"required"`
	Archive          *multipart.FileHeader `form:"archive" binding:"required"`
}

// CreateRequest создаёт заявку от пользователя.
func (h *Handler) CreateRequest(c *gin.Context) {
	var req CreateRequestReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга формы: " + err.Error()})
		return
	}

	input := CreateRequestInput{
		Title:            req.Title,
		Description:      req.Description,
		EventDate:        req.EventDate,
		EventTypeId:      req.EventTypeId,
		Email:            req.Email,
		TelegramUsername: req.TelegramUsername,
		Archive:          req.Archive,
	}

	err := h.svc.CreateRequest(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrArchiveMustBeZip):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Архив должен иметь расширение zip"})
			return
		case errors.Is(err, ErrArchiveTooLarge):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Файл должен быть не больше 50 МБ"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка создания заявки: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

// ReviewRequest берёт заявку в проверку, переводя её статус из 'Новая' в 'На проверке'.
func (h *Handler) ReviewRequest(c *gin.Context) {
	// Парсим ID заявки из URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректный ID заявки"})
		return
	}

	adminLogin := c.GetString("adminLogin")
	adminRole := c.GetString("role")

	// Берём заявку в проверку
	if err := h.svc.ReviewRequest(c.Request.Context(), id, adminLogin, adminRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ApproveRequestReq — тело запроса на одобрение заявки.
type ApproveRequestReq struct {
	Latitude    float64 `json:"latitude"    binding:"required"`
	Longitude   float64 `json:"longitude"   binding:"required"`
	EventTypeId int     `json:"eventTypeId" binding:"required,gt=0"`
}

// ApproveRequest одобряет заявку и публикует объект на карте.
func (h *Handler) ApproveRequest(c *gin.Context) {
	// Парсим ID заявки из URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректный ID заявки"})
		return
	}

	// Парсим тело запроса
	var req ApproveRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка валидации: " + err.Error()})
		return
	}

	adminLogin := c.GetString("adminLogin")
	adminRole := c.GetString("role")

	// Одобряем заявку
	err = h.svc.ApproveRequest(c.Request.Context(), id, adminLogin, adminRole, req.Latitude, req.Longitude, req.EventTypeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// RejectRequestReq — тело запроса на отклонение заявки.
type RejectRequestReq struct {
	Comment string `json:"comment" binding:"required"`
}

// RejectRequest отклоняет заявку с обязательным комментарием.
func (h *Handler) RejectRequest(c *gin.Context) {
	// Парсим ID заявки из URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректный ID заявки"})
		return
	}

	// Парсим тело запроса
	var req RejectRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Комментарий обязателен"})
		return
	}

	adminLogin := c.GetString("adminLogin")
	adminRole := c.GetString("role")

	// Отклоняем заявку
	err = h.svc.RejectRequest(c.Request.Context(), id, adminLogin, adminRole, req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
