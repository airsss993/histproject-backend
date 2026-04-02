package notifier

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Notifier — отправка вебхуков в n8n.
type Notifier struct {
	webhookURL    string
	webhookSecret string
	client        *http.Client
}

// New создаёт Notifier. Если webhookURL пустой — уведомления отключены.
func New(webhookURL, webhookSecret string) *Notifier {
	return &Notifier{
		webhookURL:    webhookURL,
		webhookSecret: webhookSecret,
		client:        &http.Client{Timeout: 5 * time.Second},
	}
}

type adminPayload struct {
	Event            string   `json:"event"`
	RequestID        int      `json:"requestId"`
	Title            string   `json:"title"`
	To               []string `json:"to"`
	Email            string   `json:"email"`
	TelegramUsername string   `json:"telegramUsername"`
}

type userPayload struct {
	Event            string  `json:"event"`
	RequestID        int     `json:"requestId"`
	Title            string  `json:"title"`
	Email            string  `json:"email"`
	TelegramUsername string  `json:"telegramUsername"`
	ChatID           *int64  `json:"chatId"`
	AdminComment     *string `json:"adminComment"`
	SiteURL          *string `json:"siteUrl"`
}

// NotifyNewRequest уведомляет администраторов о новой заявке.
func (n *Notifier) NotifyNewRequest(requestID int, title, email, telegramUsername string, adminEmails []string) {
	n.send(adminPayload{
		Event:            "new_request",
		RequestID:        requestID,
		Title:            title,
		To:               adminEmails,
		Email:            email,
		TelegramUsername: telegramUsername,
	})
}

// NotifyApproved уведомляет пользователя об одобрении заявки.
func (n *Notifier) NotifyApproved(requestID int, title, email, telegramUsername, siteURL string, chatID *int64) {
	n.send(userPayload{
		Event:            "request_approved",
		RequestID:        requestID,
		Title:            title,
		Email:            email,
		TelegramUsername: telegramUsername,
		ChatID:           chatID,
		AdminComment:     nil,
		SiteURL:          &siteURL,
	})
}

// NotifyRejected уведомляет пользователя об отклонении заявки.
func (n *Notifier) NotifyRejected(requestID int, title, email, telegramUsername, comment string, chatID *int64) {
	n.send(userPayload{
		Event:            "request_rejected",
		RequestID:        requestID,
		Title:            title,
		Email:            email,
		TelegramUsername: telegramUsername,
		ChatID:           chatID,
		AdminComment:     &comment,
		SiteURL:          nil,
	})
}

// send отправляет вебхук асинхронно; ошибки логируются как предупреждения.
func (n *Notifier) send(p any) {
	if n.webhookURL == "" {
		return
	}
	go func() {
		body, err := json.Marshal(p)
		if err != nil {
			log.Printf("[notifier] предупреждение: ошибка сериализации: %v", err)
			return
		}
		req, err := http.NewRequest(http.MethodPost, n.webhookURL, bytes.NewReader(body))
		if err != nil {
			log.Printf("[notifier] предупреждение: ошибка создания запроса: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if n.webhookSecret != "" {
			req.Header.Set("X-Webhook-Secret", n.webhookSecret)
		}
		resp, err := n.client.Do(req)
		if err != nil {
			log.Printf("[notifier] предупреждение: не удалось отправить вебхук: %v", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			log.Printf("[notifier] предупреждение: вебхук вернул статус %d", resp.StatusCode)
		}
	}()
}
