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
func (s *Service) OnNewRequest(ctx context.Context, requestID int, title, email, username string) {
	// Получаем email-адреса администраторов
	adminEmails, _ := s.repo.GetAdminEmails(ctx)
	s.webhook.NotifyNewRequest(requestID, title, email, username, adminEmails)
}

// OnUnderReview уведомляет пользователя о том, что заявка взята на проверку.
func (s *Service) OnUnderReview(ctx context.Context, requestID int, title, email, username string) {
	// Получаем chatId пользователя если он подписан на бота
	chatID, _ := s.repo.GetChatID(ctx, username)
	s.webhook.NotifyUnderReview(requestID, title, email, username, chatID)
}

// OnApproved уведомляет пользователя об одобрении заявки.
func (s *Service) OnApproved(ctx context.Context, requestID int, title, email, username, siteURL string) {
	// Получаем chatId пользователя если он подписан на бота
	chatID, _ := s.repo.GetChatID(ctx, username)
	s.webhook.NotifyApproved(requestID, title, email, username, siteURL, chatID)
}

// OnRejected уведомляет пользователя об отклонении заявки.
func (s *Service) OnRejected(ctx context.Context, requestID int, title, email, username, comment string) {
	// Получаем chatId пользователя если он подписан на бота
	chatID, _ := s.repo.GetChatID(ctx, username)
	s.webhook.NotifyRejected(requestID, title, email, username, comment, chatID)
}

// OnAdminCreated уведомляет нового администратора о его учётных данных.
func (s *Service) OnAdminCreated(_ context.Context, email, login, password, role string) {
	s.webhook.NotifyAdminCreated(email, login, password, role)
}
