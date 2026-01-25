package main

import (
	_ "github.com/airsss993/histproject-backend/docs"
	"github.com/airsss993/histproject-backend/internal/app"
)

// @title Документация API исторической платформы
// @version 1.0
// @description Серверное REST API исторической платформы, разработанное на Go с использованием Gin.
// @host localhost:7666
// @BasePath /api
func main() {
	app.Run()
}
