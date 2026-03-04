package router

import (
	"net/http"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/internal/objects"
	"github.com/airsss993/histproject-backend/internal/requests"
	"github.com/gin-gonic/gin"
)

// Handlers — все HTTP-хендлеры приложения.
type Handlers struct {
	Objects  *objects.Handler
	Requests *requests.Handler
}

func New(cfg *config.Config, h Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware(cfg.CORS.AllowedOrigins))

	InitRoutes(r, h)

	return r
}

func InitRoutes(r *gin.Engine, h Handlers) {
	public := r.Group("/api")
	{
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		public.GET("objects/get-object-data/:id", h.Objects.GetObjectData)
		public.POST("objects/get-objects-list", h.Objects.GetObjectsList)
		public.GET("objects/get-event-types-list", h.Objects.GetEventTypesList)

		public.POST("requests/create-request", h.Requests.CreateRequest)
	}
}

func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//origin := c.Request.Header.Get("Origin")
		//
		//if origin != "" && strings.Contains(allowedOrigins, origin) {
		//	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
