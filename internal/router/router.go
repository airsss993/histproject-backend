package router

import (
	"net/http"

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
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
}
