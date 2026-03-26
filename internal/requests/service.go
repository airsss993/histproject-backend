package requests

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/hibiken/asynq"
)

const maxArchiveSize = 50 * 1024 * 1024 // 50 МБ

// StorageWriter — загрузка архива в хранилище.
type StorageWriter interface {
	UploadArchive(file *multipart.FileHeader) (string, error)
}

// QueueEnqueuer — постановка задачи в очередь.
type QueueEnqueuer interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

// TaskBuilder создаёт задачу обработки архива.
type TaskBuilder func(requestID int, archiveID string) (*asynq.Task, error)

// Service — бизнес-логика заявок.
type Service struct {
	repo        RequestRepository
	storage     StorageWriter
	queue       QueueEnqueuer
	taskBuilder TaskBuilder
}

// NewService создаёт сервис заявок.
func NewService(repo RequestRepository, storage StorageWriter, queue QueueEnqueuer, taskBuilder TaskBuilder) *Service {
	return &Service{
		repo:        repo,
		storage:     storage,
		queue:       queue,
		taskBuilder: taskBuilder,
	}
}

// CreateRequestInput — входные данные для создания заявки.
type CreateRequestInput struct {
	Title            string
	Description      string
	EventDate        string
	EventTypeId      int
	Email            string
	TelegramUsername string
	Archive          *multipart.FileHeader
}

// CreateRequest создаёт заявку: проверяет архив, загружает, сохраняет в БД, ставит задачу в очередь.
func (s *Service) CreateRequest(ctx context.Context, input CreateRequestInput) (err error) {
	if !strings.Contains(input.Archive.Filename, ".zip") {
		return ErrArchiveMustBeZip
	}
	if input.Archive.Size > maxArchiveSize {
		return ErrArchiveTooLarge
	}

	archiveId, err := s.storage.UploadArchive(input.Archive)
	if err != nil {
		return err
	}

	data := RequestData{
		Title:            input.Title,
		Description:      input.Description,
		Email:            input.Email,
		TelegramUsername: input.TelegramUsername,
		EventDate:        input.EventDate,
		EventTypeId:      input.EventTypeId,
		ArchiveId:        archiveId,
	}

	requestID, err := s.repo.Create(ctx, data)
	if err != nil {
		return err
	}

	task, err := s.taskBuilder(requestID, archiveId)
	if err != nil {
		return err
	}

	_, _ = s.queue.Enqueue(task)
	return nil
}

