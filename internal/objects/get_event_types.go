package objects

import (
	"fmt"
	"net/http"

	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetEventTypesListResp - структура ответа для получения списка типов событий
type GetEventTypesListResp struct {
	// EventTypes - массив словарей, каждый элемент которого содержит:
	// id - ID типа события
	// name - название типа события
	// description - описание типа события
	EventTypes []EventTypeInfo `json:"eventTypes"`
}

// @Summary		Получение списка типов событий
// @Description	Метод для получения списка типов событий.
// @Tags		Метки
// @Produce		json
// @Success		200	{events} EventTypeInfo
// @Router		/objects/get-event-types-list [get]
func GetEventTypesList(c *gin.Context) {
	// Инициализируем структуру ответа
	resp := GetEventTypesListResp{}

	// Получаем список типов событий и их данные
	eventTypes, err := getEventTypesListFromDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка получения данных из БД: " + err.Error(),
		})
		return
	}

	// Отправляем успешный ответ
	resp.EventTypes = eventTypes
	c.JSON(http.StatusOK, resp)
}

// Модель объекта (как он лежит в БД и как уходит в JSON)
type EventTypeInfo struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

// getEventTypesListFromDB - функция для получения всех типов событий с их данными из БД
func getEventTypesListFromDB() ([]EventTypeInfo, error) {
	var eventTypesInfo []EventTypeInfo

	query := `
		SELECT
		    id,
		    name,
		    description
		FROM
		    histproject.event_types`

	err := db.DB.Select(&eventTypesInfo, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения типов событий с БД: %w", err)
	}

	return eventTypesInfo, nil
}
