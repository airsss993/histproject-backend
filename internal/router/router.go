package router

import (
	"net/http"

	_ "github.com/airsss993/histproject-backend/docs"
	"github.com/airsss993/histproject-backend/internal/objects"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.Default()

	InitRoutes(r)

	return r
}

func InitRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Публичные роуты, не требующие проверки авторизации
	public := r.Group("/api")
	{
		// Тестовый эндпоинт для проверки работы сервера
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// ---------- ОБЪЕКТЫ И ТИПЫ СОБЫТИЙ ----------

		// Эндпоинт для получения информации об одном объекте
		public.GET("objects/get-object-data/:id", objects.GetObjectData)
		// Эндпоинт для получения всех объектов с их информацией
		public.GET("objects/get-objects-list", objects.GetObjectsList)
		// Эндпоинт для получения всех типов событий с их информацией
		public.GET("objects/get-event-types-list", objects.GetEventTypesList)
	}
}
