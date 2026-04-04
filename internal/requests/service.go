package requests

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/airsss993/histproject-backend/internal/notifications"
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
	repo          RequestRepository
	storage       StorageWriter
	queue         QueueEnqueuer
	taskBuilder   TaskBuilder
	objectsRepo   objects.Repository
	notifications *notifications.Service
}

// NewService создаёт сервис заявок.
func NewService(repo RequestRepository, storage StorageWriter, queue QueueEnqueuer, taskBuilder TaskBuilder, objectsRepo objects.Repository, notifications *notifications.Service) *Service {
	return &Service{
		repo:          repo,
		storage:       storage,
		queue:         queue,
		taskBuilder:   taskBuilder,
		objectsRepo:   objectsRepo,
		notifications: notifications,
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
	// Проверяем, что файл — zip и не больше 50 МБ
	if !strings.Contains(input.Archive.Filename, ".zip") {
		return ErrArchiveMustBeZip
	}
	if input.Archive.Size > maxArchiveSize {
		return ErrArchiveTooLarge
	}

	// Загружаем архив в хранилище и получаем archiveId
	archiveId, err := s.storage.UploadArchive(input.Archive)
	if err != nil {
		return err
	}

	// Сохраняем заявку в БД со статусом 'В обработке'
	data := RequestData{
		Title:            input.Title,
		Description:      input.Description,
		Email:            input.Email,
		TelegramUsername: input.TelegramUsername,
		EventDate:        input.EventDate,
		EventTypeId:      input.EventTypeId,
		ArchiveId:        archiveId,
	}

	// Получаем ID новой заявки
	requestID, err := s.repo.Create(ctx, data)
	if err != nil {
		return err
	}

	// Ставим задачу в очередь на обработку архива
	task, err := s.taskBuilder(requestID, archiveId)
	if err != nil {
		return err
	}

	_, _ = s.queue.Enqueue(task)

	// Уведомляем администраторов о новой заявке
	s.notifications.OnNewRequest(ctx, requestID, input.Title, input.Email, input.TelegramUsername)
	return nil
}

// ReviewRequest переводит заявку в статус 'На проверке' для дальнейшей модерации.
func (s *Service) ReviewRequest(ctx context.Context, id int, adminLogin, adminRole string) error {
	// Получаем заявку из БД
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("заявка не найдена")
	}

	// Брать в проверку можно только заявки со статусом 'Новая'
	if req.Status != "Новая" {
		return errors.New("заявка не находится в статусе 'Новая'")
	}

	// Переводим статус заявки в 'На проверке'
	if err := s.repo.UpdateStatus(ctx, id, "На проверке", "", req.SiteURL, req.ScreenshotURL); err != nil {
		return err
	}

	// Записываем действие в аудит-лог
	_ = s.repo.LogHistory(ctx, id, adminLogin, adminRole, "Взята на проверку", "")

	return nil
}

// ApproveRequest одобряет заявку: создаёт объект на карте и переводит статус в 'Опубликована'.
func (s *Service) ApproveRequest(ctx context.Context, id int, adminLogin, adminRole string, latitude, longitude float64, eventTypeID int) error {
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
	if err := s.repo.UpdateStatus(ctx, id, "Опубликована", "", req.SiteURL, req.ScreenshotURL); err != nil {
		return err
	}

	// Записываем действие в аудит-лог
	_ = s.repo.LogHistory(ctx, id, adminLogin, adminRole, "Одобрена", "")

	// Уведомляем пользователя об одобрении заявки
	s.notifications.OnApproved(ctx, req.ID, req.Title, req.Email, req.TelegramUsername, req.SiteURL)
	return nil
}

// RejectRequest отклоняет заявку с обязательным комментарием.
func (s *Service) RejectRequest(ctx context.Context, id int, adminLogin, adminRole, comment string) error {
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
	if err := s.repo.UpdateStatus(ctx, id, "Отклонена", comment, "", ""); err != nil {
		return err
	}

	// Записываем действие в аудит-лог
	_ = s.repo.LogHistory(ctx, id, adminLogin, adminRole, "Отклонена", comment)

	// Уведомляем пользователя об отклонении заявки
	s.notifications.OnRejected(ctx, req.ID, req.Title, req.Email, req.TelegramUsername, comment)
	return nil
}
