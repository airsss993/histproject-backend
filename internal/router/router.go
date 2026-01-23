package router

import (
	"net/http"

	"github.com/airsss993/histproject-backend/internal/objects"
	"github.com/gin-gonic/gin"
)

func New(basePath string) *gin.Engine {
	r := gin.Default()

	InitRoutes(r, basePath)

	return r
}

func InitRoutes(r *gin.Engine, basePath string) {
	// Публичные роуты, не требующие проверки авторизации
	public := r.Group(basePath)
	{
		// Тестовый эндпоинт для проверки работы сервера
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// ---------- ОБЪЕКТЫ И ТИПЫ СОБЫТИЙ ----------

		// Эндпоинт для получения информации об одном объекте
		public.GET("objects/get-object-details/:id", objects.GetObjectDetails)
		// Эндпоинт для получения всех объектов с их информацией
		public.GET("objects/get-objects-list", objects.GetObjectsList)
		// Эндпоинт для получения всех типов событий с их информацией
		public.GET("objects/get-event-types-list", objects.GetEventTypesList)
	}
}
