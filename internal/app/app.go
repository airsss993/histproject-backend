package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/internal/router"
	"github.com/airsss993/histproject-backend/internal/server"
	"github.com/airsss993/histproject-backend/internal/worker"
	"github.com/airsss993/histproject-backend/migrations"
	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/airsss993/histproject-backend/pkg/queue"
	"github.com/airsss993/histproject-backend/pkg/storage"
)

func Run() {
	// Инициализируем конфиг приложения
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига: ", err)
	}

	// Создаем подключение к БД
	conn := db.ConnDB(cfg.Database.DSN)

	minioClient := storage.NewMinioClient(cfg.Storage)
	if err := minioClient.InitMinio(); err != nil {
		log.Fatal("Ошибка подключения к MinIO: ", err)
	}
	// Устанавливаем глобальный экземпляр MinIO клиента
	storage.Client = minioClient

	queue.QueuClient = queue.NewClient(cfg.Redis)
	queueServer := queue.NewServer(cfg.Redis)

	w := worker.NewWorker(conn, minioClient, *cfg)
	mux := worker.NewMux(w)
	go queueServer.Run(mux)

	// Выполняем миграции
	if err := migrations.Run(conn.DB); err != nil {
		log.Fatal("Ошибка выполнения миграций: ", err)
	}

	// Создаем роутер
	r := router.New(cfg)

	// Создаем сервер и запускаем его
	srv := server.New(cfg.App.Port, r)
	srv.Start()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Завершение работы сервера...")

	// Даём серверу 5 секунд на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем сервер
	srv.Stop(ctx)
}
