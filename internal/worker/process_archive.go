package worker

import (
	"context"
	"fmt"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/pkg/storage"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

// Worker - структура для обработки задач в фоновом режиме
type Worker struct {
	db      *sqlx.DB
	storage *storage.MinioClient
	cfg     config.Config // TODO: передавать конфиг для воркера
}

// NewWorker - функция для создания нового экземпляра Worker
func NewWorker(db *sqlx.DB, storage *storage.MinioClient, cfg config.Config) *Worker {
	return &Worker{
		db:      db,
		storage: storage,
		cfg:     cfg,
	}
}

// NewMux - функция для создания нового ServeMux и регистрации обработчиков задач
func NewMux(w *Worker) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TypeProcessArchive, w.ProcessArchiveTask)

	return mux
}

// ProcessArchiveTask - обработчик задачи по обработке архива
func (w *Worker) ProcessArchiveTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println("process archive task")

	return nil
}
