package requests

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/airsss993/histproject-backend/internal/objects"
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
	objectsRepo objects.Repository
}

// NewService создаёт сервис заявок.
func NewService(repo RequestRepository, storage StorageWriter, queue QueueEnqueuer, taskBuilder TaskBuilder, objectsRepo objects.Repository) *Service {
	return &Service{
		repo:        repo,
		storage:     storage,
		queue:       queue,
		taskBuilder: taskBuilder,
		objectsRepo: objectsRepo,
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

// ApproveRequest одобряет заявку: создаёт объект на карте и переводит статус в 'Опубликована'.
func (s *Service) ApproveRequest(ctx context.Context, id int, latitude, longitude float64, eventTypeID int) error {
	// Получаем заявку из БД
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("заявка не найдена")
	}

	// Одобрять можно только заявки со статусом 'На проверке'
	if req.Status != "На проверке" {
		return errors.New("заявка не находится в статусе 'На проверке'")
	}

	// Создаём объект на карте
	if err := s.objectsRepo.CreateObject(ctx, objects.ObjectData{
		RequestID:       req.ID,
		Title:           req.Title,
		Description:     req.Description,
		Latitude:        latitude,
		Longitude:       longitude,
		EventDate:       req.EventDate,
		EventTypeID:     eventTypeID,
		SiteURL:         req.SiteURL,
		PreviewImageURL: req.ScreenshotURL,
	}); err != nil {
		return err
	}

	// Переводим статус заявки в 'Опубликована'
	return s.repo.UpdateStatus(ctx, id, "Опубликована", "", req.SiteURL, req.ScreenshotURL)
}

// RejectRequest отклоняет заявку с обязательным комментарием.
func (s *Service) RejectRequest(ctx context.Context, id int, comment string) error {
	// Проверяем наличие комментария
	if strings.TrimSpace(comment) == "" {
		return errors.New("комментарий обязателен при отклонении")
	}

	// Получаем заявку из БД
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("заявка не найдена")
	}

	// Отклонять можно только заявки со статусом 'На проверке'
	if req.Status != "На проверке" {
		return errors.New("заявка не находится в статусе 'На проверке'")
	}

	// Переводим статус заявки в 'Отклонена'
	return s.repo.UpdateStatus(ctx, id, "Отклонена", comment, "", "")
}
