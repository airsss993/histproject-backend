package router

import (
	"net/http"
	"strings"

	"github.com/airsss993/histproject-backend/docs"
	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/internal/objects"
	"github.com/airsss993/histproject-backend/internal/requests"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware(cfg.CORS.AllowedOrigins))

	InitRoutes(r, cfg.App.SwaggerHost)

	return r
}

func InitRoutes(r *gin.Engine, swaggerHost string) {
	docs.SwaggerInfo.Host = swaggerHost
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
		public.POST("objects/get-objects-list", objects.GetObjectsList)
		// Эндпоинт для получения всех типов событий с их информацией
		public.GET("objects/get-event-types-list", objects.GetEventTypesList)

		// ---------- ПОЛЬЗОВАТЕЛЬСКИЕ ЗАЯВКИ ----------

		// Эндпоинт для создания заявки от пользователя
		public.POST("requests/create-request", requests.CreateRequest)
	}
}

func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" && strings.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
