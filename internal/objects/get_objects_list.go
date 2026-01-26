package objects

import (
	"fmt"
	"net/http"

	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// GetObjectsListReq - структура запроса для получения списка объектов с фильтрацией
type GetObjectsListReq struct {
	EventTypeIDs []int  `json:"eventTypeIds" binding:"omitempty,dive,gt=0"`
	DateFrom     string `json:"dateFrom" binding:"omitempty,datetime=2006-01-02"`
	DateTo       string `json:"dateTo" binding:"omitempty,datetime=2006-01-02"`
}

// GetObjectsListResp - структура ответа для получения списка объектов
type GetObjectsListResp struct {
	// Objects - массив словарей, каждый элемент которого содержит:
	// id - ID объекта
	// title - название объекта
	// description - описание объекта
	// coordinates - координаты объекта
	// eventDate - дата события объекта
	// eventType - тип объекта
	// previewUrlImage - URL картинки превью объекта
	Objects []ObjectInfo `json:"objects"`
}

// @Summary		Получение списка объектов
// @Description	Метод для получения списка объектов.
// @Tags		Объекты
// @Produce		json
// @Success		200	{objects} ObjectInfo
// @Router		/objects/get-objects-list [get]
func GetObjectsList(c *gin.Context) {
	// Инициализируем структуру запроса для того чтобы заполнить её данными
	var req GetObjectsListReq

	if c.Request.ContentLength > 0 {
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Ошибка парсинга JSON: " + err.Error(),
			})
			return
		}
	}

	// Получаем список объектов с фильтрами
	objects, err := getObjectsListFromDB(req.EventTypeIDs, req.DateFrom, req.DateTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка получения данных из БД: " + err.Error(),
		})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, GetObjectsListResp{Objects: objects})

}

// getObjectsListFromDB - функция для получения информации о всех объектах из БД.
func getObjectsListFromDB(eventTypesIds []int, dateFrom, dateTo string) ([]ObjectInfo, error) {
	var objectsInfo []ObjectInfo

	query := `
		SELECT
		    id,
		    title,
		    description,
		    event_date,
		    event_type_id,
		    preview_image_url
		FROM
		    histproject.objects`

	var args []interface{}
	argPos := 1

	// Если пришли типы, то добавляем фильтрацию по ним
	if len(eventTypesIds) > 0 {
		query += fmt.Sprintf(" WHERE event_type_id = ANY($%d)", argPos)
		args = append(args, pq.Array(eventTypesIds))
		argPos++
	}

	// Если пришла дата начала, то добавляем
	if dateFrom != "" {
		if argPos > 1 {
			query += fmt.Sprintf(" AND event_date >= $%d", argPos)
		} else {
			query += fmt.Sprintf(" WHERE event_date >= $%d", argPos)
		}
		args = append(args, dateFrom)
		argPos++
	}

	// Если пришла дата конца, то добавляем
	if dateTo != "" {
		if argPos > 1 {
			query += fmt.Sprintf(" AND event_date <= $%d", argPos)
		} else {
			query += fmt.Sprintf(" WHERE event_date <= $%d", argPos)
		}
		args = append(args, dateTo)
		argPos++
	}

	err := db.DB.Select(&objectsInfo, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объектов с БД: %w", err)
	}

	if objectsInfo == nil {
		return []ObjectInfo{}, nil
	}

	return objectsInfo, nil
}
