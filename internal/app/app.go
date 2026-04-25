package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/airsss993/histproject-backend/internal/admin"
	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/internal/notifications"
	"github.com/airsss993/histproject-backend/internal/objects"
	"github.com/airsss993/histproject-backend/internal/requests"
	"github.com/airsss993/histproject-backend/internal/router"
	"github.com/airsss993/histproject-backend/internal/server"
	"github.com/airsss993/histproject-backend/internal/worker"
	"github.com/airsss993/histproject-backend/migrations"
	"github.com/airsss993/histproject-backend/pkg/db"
	"github.com/airsss993/histproject-backend/pkg/notifier"
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
	objectsSvc := objects.NewService(objectsRepo, minioClient)
	objectsHandler := objects.NewHandler(objectsSvc)

	// Инициализация модуля admin
	adminRepo := admin.NewRepository(conn)
	adminSvc := admin.NewService(adminRepo, requestsRepo, cfg.Auth.JWTSecret)

	// Инициализация модуля notifications
	notifRepo := notifications.NewRepository(conn)
	n8nNotifier := notifier.New(cfg.N8N.WebhookURL, cfg.N8N.WebhookSecret)
	notifSvc := notifications.New(n8nNotifier, notifRepo)

	requestsSvc := requests.NewService(requestsRepo, minioClient, queueClient, worker.NewProcessArchiveTask, objectsRepo, notifSvc)
	requestsHandler := requests.NewHandler(requestsSvc)
	adminHandler := admin.NewHandler(adminSvc)

	// Автосоздание super_admin при первом запуске
	if err := adminSvc.SeedSuperAdmin(context.Background()); err != nil {
		log.Fatal("Ошибка создания super_admin: ", err)
	}

	r := router.New(cfg, router.Handlers{
		Objects:  objectsHandler,
		Requests: requestsHandler,
		Admin:    adminHandler,
	}, adminSvc)

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
