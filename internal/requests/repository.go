package requests

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// RequestRepository — доступ к заявкам (создание и обновление статуса).
type RequestRepository interface {
	Create(ctx context.Context, data RequestData) (int, error)
	UpdateStatus(ctx context.Context, id int, status, comment, siteURL, screenshotURL string) error
	List(ctx context.Context, status string) ([]RequestListItem, error)
	GetByID(ctx context.Context, id int) (*RequestDetail, error)
}

// RequestData — данные заявки для сохранения в БД.
type RequestData struct {
	Title            string
	Description      string
	Email            string
	TelegramUsername string
	ArchiveId        string
	SiteURL          string
	ScreenshotURL    string
	EventDate        string
	EventTypeId      int
}

type repository struct {
	db *sqlx.DB
}

// NewRepository создаёт реализацию RequestRepository для PostgreSQL.
func NewRepository(db *sqlx.DB) RequestRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data RequestData) (int, error) {
	var id int
	query := `
		INSERT INTO requests 
			(title, description, email, telegram_username, archive_id, site_url, screenshot_url, event_date, event_type_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		data.Title, data.Description, data.Email, data.TelegramUsername, data.ArchiveId,
		data.SiteURL, data.ScreenshotURL, data.EventDate, data.EventTypeId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки новой записи о заявке в БД: %w", err)
	}
	return id, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id int, status, comment, siteURL, screenshotURL string) error {
	query := `UPDATE requests SET status = $1, admin_comment = $2, site_url = $3, screenshot_url = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, status, comment, siteURL, screenshotURL, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса заявки: %w", err)
	}
	return nil
}

func (r *repository) List(ctx context.Context, status string) ([]RequestListItem, error) {
	query := `SELECT id, title, status, created_at FROM requests`
	var args []interface{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC`

	var list []RequestListItem
	err := r.db.SelectContext(ctx, &list, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка заявок: %w", err)
	}
	if list == nil {
		return []RequestListItem{}, nil
	}
	return list, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*RequestDetail, error) {
	var req RequestDetail
	query := `
		SELECT id, title, description, event_date, event_type_id, email, telegram_username,
		       status, admin_comment, site_url, screenshot_url, created_at
		FROM requests WHERE id = $1`
	err := r.db.GetContext(ctx, &req, query, id)
	if err != nil {
		return nil, fmt.Errorf("заявка не найдена: %w", err)
	}
	return &req, nil
}
