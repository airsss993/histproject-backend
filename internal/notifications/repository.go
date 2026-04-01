package notifications

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository — чтение подписок на Telegram-бота.
type Repository struct {
	db *sqlx.DB
}

// NewRepository создаёт Repository.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// GetAdminEmails возвращает список email-адресов всех администраторов.
func (r *Repository) GetAdminEmails(ctx context.Context) ([]string, error) {
	var emails []string
	query := `SELECT email FROM admins WHERE email != ''`
	err := r.db.SelectContext(ctx, &emails, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения email-адресов администраторов: %w", err)
	}
	return emails, nil
}

// GetChatID возвращает chatId по telegramUsername. Если подписки нет — возвращает nil.
func (r *Repository) GetChatID(ctx context.Context, username string) (*int64, error) {
	var chatID int64
	query := `SELECT chat_id FROM telegram_subscriptions WHERE telegram_username = $1`
	err := r.db.GetContext(ctx, &chatID, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения chatId: %w", err)
	}
	return &chatID, nil
}
