package objects

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler — HTTP-обработчики для объектов и типов событий.
type Handler struct {
	svc *Service
}

// NewHandler создаёт handler объектов.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GetObjectDataReq — тело/параметры запроса получения данных объекта.
type GetObjectDataReq struct {
	ObjectId int `json:"objectId" binding:"required,gt=0"`
}

// GetObjectDataResp — ответ с данными одного объекта.
type GetObjectDataResp struct {
	Object SingleObjectInfo `json:"object"`
}

// GetObjectData возвращает данные одного объекта по ID.
// @Summary		Получение данных для конкретного объекта
// @Description	Метод для получения данных о конкретном объекте (метки на карте)
// @Tags		Объекты
// @Produce		json
// @Param		id	path		int	true	"ID получаемого объекта (метки на карте)"
// @Success		200	{object}	map[string]ObjectInfo
// @Failure		400	{object}	map[string]string
// @Router		/objects/get-object-data/{id} [get]
func (h *Handler) GetObjectData(c *gin.Context) {
	objectIdParam := c.Param("id")
	objectId, err := strconv.Atoi(objectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка валидации: ID объекта должен быть числом"})
		return
	}

	obj, err := h.svc.GetObjectData(c.Request.Context(), objectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка получения данных из БД: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"object": *obj})
}

// GetObjectsListReq — параметры фильтрации списка объектов.
type GetObjectsListReq struct {
	EventTypeIDs []int  `json:"eventTypeIds" binding:"omitempty,dive,gt=0"`
	DateFrom     string `json:"dateFrom" binding:"omitempty,datetime=2006-01-02"`
	DateTo       string `json:"dateTo" binding:"omitempty,datetime=2006-01-02"`
}

// GetObjectsListResp — ответ со списком объектов.
type GetObjectsListResp struct {
	Objects []ObjectInfo `json:"objects"`
}

// GetObjectsList возвращает список объектов с фильтрами.
// @Summary		Получение списка объектов
// @Description	Метод для получения списка объектов с возможностью фильтрации по типам событий и датам
// @Tags		Объекты
// @Accept		json
// @Produce		json
// @Param		request	body	GetObjectsListReq	false	"Параметры фильтрации (опционально)"
// @Success		200	{object}	GetObjectsListResp
// @Failure		400	{object}	map[string]string
// @Router		/objects/get-objects-list [post]
func (h *Handler) GetObjectsList(c *gin.Context) {
	var req GetObjectsListReq
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга JSON: " + err.Error()})
			return
		}
	}

	list, err := h.svc.GetObjectsList(c.Request.Context(), req.EventTypeIDs, req.DateFrom, req.DateTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка получения данных из БД: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetObjectsListResp{Objects: list})
}

// GetEventTypesListResp — ответ со списком типов событий.
type GetEventTypesListResp struct {
	EventTypes []EventTypeInfo `json:"eventTypes"`
}

// GetEventTypesList возвращает список типов событий.
// @Summary		Получение списка типов событий
// @Description	Метод для получения списка типов событий
// @Tags		Объекты
// @Produce		json
// @Success		200	{object}	GetEventTypesListResp
// @Failure		400	{object}	map[string]string
// @Router		/objects/get-event-types-list [get]
func (h *Handler) GetEventTypesList(c *gin.Context) {
	list, err := h.svc.GetEventTypesList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка получения данных из БД: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, GetEventTypesListResp{EventTypes: list})
}
