package notifications

import (
	"context"

	"github.com/airsss993/histproject-backend/pkg/notifier"
)

// Service — оркестрирует отправку уведомлений о событиях с заявками.
type Service struct {
	webhook *notifier.Notifier
	repo    *Repository
}

// New создаёт сервис уведомлений.
func New(webhook *notifier.Notifier, repo *Repository) *Service {
	return &Service{webhook: webhook, repo: repo}
}

// OnNewRequest уведомляет администраторов о новой заявке.
func (s *Service) OnNewRequest(ctx context.Context, requestID int, title, email, telegramUsername string) {
	// Получаем email-адреса администраторов
	adminEmails, _ := s.repo.GetAdminEmails(ctx)
	s.webhook.NotifyNewRequest(requestID, title, email, telegramUsername, adminEmails)
}

// OnApproved уведомляет пользователя об одобрении заявки.
func (s *Service) OnApproved(ctx context.Context, requestID int, title, email, telegramUsername, siteURL string) {
	// Получаем chatId пользователя если он подписан на бота
	chatID, _ := s.repo.GetChatID(ctx, telegramUsername)
	s.webhook.NotifyApproved(requestID, title, email, telegramUsername, siteURL, chatID)
}

// OnRejected уведомляет пользователя об отклонении заявки.
func (s *Service) OnRejected(ctx context.Context, requestID int, title, email, telegramUsername, comment string) {
	// Получаем chatId пользователя если он подписан на бота
	chatID, _ := s.repo.GetChatID(ctx, telegramUsername)
	s.webhook.NotifyRejected(requestID, title, email, telegramUsername, comment, chatID)
}
