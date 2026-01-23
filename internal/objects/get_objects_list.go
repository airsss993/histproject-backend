package objects

import (
	"fmt"
	"net/http"

	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/gin-gonic/gin"
)

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

// GetObjectsList - HTTP-хендлер для получения списка объектов (меток на карте)
func GetObjectsList(c *gin.Context) {
	// Инициализируем структуру ответа
	//out := GetObjectsListResp{
	//	Objects: []map[string]string{},
	//}

	resp := GetObjectsListResp{}

	// Получаем данные об одном объекте по его ID
	objects, err := getObjectsListFromDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка получения данных из БД: " + err.Error(),
		})
		return
	}

	resp.Objects = objects

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"objects": resp.Objects,
	})

}

// getObjectsListFromDB - функция для получения информации о всех объектах из БД.
func getObjectsListFromDB() ([]ObjectInfo, error) {
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

	err := db.DB.Select(&objectsInfo, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объектов с БД: %w", err)
	}

	return objectsInfo, nil
}
