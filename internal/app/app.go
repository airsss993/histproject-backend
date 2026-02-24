package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/internal/objects"
	"github.com/airsss993/histproject-backend/internal/requests"
	"github.com/airsss993/histproject-backend/internal/router"
	"github.com/airsss993/histproject-backend/internal/server"
	"github.com/airsss993/histproject-backend/internal/worker"
	"github.com/airsss993/histproject-backend/migrations"
	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/airsss993/histproject-backend/pkg/queue"
	"github.com/airsss993/histproject-backend/pkg/storage"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига: ", err)
	}

	conn := db.ConnDB(cfg.Database.DSN)

	minioClient := storage.NewMinioClient(cfg.Storage)
	if err := minioClient.InitMinio(); err != nil {
		log.Fatal("Ошибка подключения к MinIO: ", err)
	}
	storage.Client = minioClient

	queueClient := queue.NewClient(cfg.Redis)
	queue.QueuClient = queueClient
	queueServer := queue.NewServer(cfg.Redis)

	requestsRepo := requests.NewRepository(conn)
	w := worker.NewWorker(requestsRepo, minioClient, *cfg)
	mux := worker.NewMux(w)
	go queueServer.Run(mux)

	if err := migrations.Run(conn.DB); err != nil {
		log.Fatal("Ошибка выполнения миграций: ", err)
	}

	objectsRepo := objects.NewRepository(conn)
	objectsSvc := objects.NewService(objectsRepo)
	objectsHandler := objects.NewHandler(objectsSvc)

	requestsSvc := requests.NewService(requestsRepo, minioClient, queueClient, worker.NewProcessArchiveTask)
	requestsHandler := requests.NewHandler(requestsSvc)

	r := router.New(cfg, router.Handlers{
		Objects:  objectsHandler,
		Requests: requestsHandler,
	})

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
