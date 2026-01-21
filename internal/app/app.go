package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/internal/router"
	"github.com/airsss993/histproject-backend/pkg/db"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига: ", err)
	}
	_ = db.ConnDB(cfg)

	r := router.New(cfg.Server.BasePath)

	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Println("[INFO] Сервер запущен на порту", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("[ERROR] Ошибка запуска сервера: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("[ERROR] Ошибка при завершении сервера: ", err)
	}

	log.Println("[INFO] Сервер остановлен")
}
