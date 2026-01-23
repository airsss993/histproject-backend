package objects

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetObjectDetailsReq - структура запроса для получения данных о конкретном объекта
type GetObjectDetailsReq struct {
	ObjectId int `json:"objectId" bindind:"required,gt=0"` // ID объекта (метки на карте), для которого надо получить данные
}

// GetObjectDetailsResp - структура ответа для получения данных о конкретном объекта
type GetObjectDetailsResp struct {
	// Object - словарь, который содержит информацию об одном объекте:
	// id - ID объекта
	// title - название объекта
	// description - описание объекта
	// eventDate - дата события объекта
	// eventTypeId - ID типа объекта
	// previewUrlImage - URL картинки превью объекта
	Object ObjectInfo `json:"object"`
}

// GetObjectDetails - HTTP-хендлер для получения данных для конкретного объекта
func GetObjectDetails(c *gin.Context) {
	// Инициализируем структуру запроса и ответа
	//req := GetObjectDetailsReq{
	//	ObjectId: 0,
	//}
	resp := GetObjectDetailsResp{
		Object: ObjectInfo{
			ID:              0,
			Title:           "",
			Description:     "",
			EventDate:       "",
			EventTypeID:     0,
			PreviewUrlImage: "",
		},
	}

	// Получаем ID объекта из параметров запроса
	objectIdParam := c.Param("id")
	objectId, err := strconv.Atoi(objectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка валидации: ID объекта должен быть не пустой строкой" + err.Error(),
		})
		return
	}

	// Получаем данные об одном объекте по его ID
	objectInfo, err := getObjectDetails(objectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка получения данных из БД: " + err.Error(),
		})
		return
	}

	resp.Object = *objectInfo

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"object":  resp.Object,
	})
}

// Модель объекта (как он лежит в БД и как уходит в JSON)
type ObjectInfo struct {
	ID              int    `db:"id" json:"id"`
	Title           string `db:"title" json:"title"`
	Description     string `db:"description" json:"description"`
	EventDate       string `db:"event_date" json:"eventDate"`
	EventTypeID     int    `db:"event_type_id" json:"eventType"`
	PreviewUrlImage string `db:"preview_image_url" json:"previewUrlImage"`
}

// getObjectDetails - функция для получения информации об одном объекте из БД по его ID.
func getObjectDetails(objectId int) (*ObjectInfo, error) {
	var objectInfo ObjectInfo

	query := `
		SELECT
		    id,
		    title,
		    description,
		    event_date,
		    event_type_id,
		    preview_image_url
		FROM
		    histproject.objects
		WHERE id = $1`

	err := db.DB.Get(&objectInfo, query, objectId)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации об объекте %v с БД: %w", objectId, err)
	}

	return &objectInfo, nil
}
