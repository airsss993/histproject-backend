package objects

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetObjectDataReq - структура запроса для получения данных о конкретном объекта
type GetObjectDataReq struct {
	ObjectId int `json:"objectId" bindind:"required,gt=0"` // ID объекта (метки на карте), для которого надо получить данные
}

// GetObjectDataResp - структура ответа для получения данных о конкретном объекта
type GetObjectDataResp struct {
	// Object - словарь, который содержит информацию об одном объекте:
	// id - ID объекта
	// title - название объекта
	// description - описание объекта
	// eventDate - дата события объекта
	// eventTypeId - ID типа объекта
	// previewUrlImage - URL картинки превью объекта
	Object ObjectInfo `json:"object"`
}

// @Summary		Получение данных для конкретного объекта
// @Description	Метод для получения данных о конкретном объекте (метки на карте)
// @Tags		Объекты
// @Produce		json
// @Param		id	path		int	true	"ID получаемого объекта (метки на карте)"
// @Success		200	{object} ObjectInfo
// @Router		/objects/get-object-data/{id} [get]
func GetObjectData(c *gin.Context) {
	// Инициализируем структуру запроса и ответа
	resp := GetObjectDataResp{
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
	objectInfo, err := getObjectData(objectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка получения данных из БД: " + err.Error(),
		})
		return
	}

	resp.Object = *objectInfo

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"object": resp.Object,
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

// GetObjectData- функция для получения информации об одном объекте из БД по его ID.
func getObjectData(objectId int) (*ObjectInfo, error) {
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
		    objects
		WHERE id = $1`

	err := db.DB.Get(&objectInfo, query, objectId)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации об объекте %v с БД: %w", objectId, err)
	}

	return &objectInfo, nil
}
