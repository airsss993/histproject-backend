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
	// Публичные роуты
	public := r.Group(basePath)
	{
		// Тестовый эндпоинт для проверки работы сервера
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Эндпоинт для получения информации об одном объекте
		public.GET("objects/getObjectDetails/:id", objects.GetObjectDetails)
		// Эндпоинт для получения всех объектов с их информацией
		public.GET("objects/getObjectsList", objects.GetObjectsList)
	}
}
