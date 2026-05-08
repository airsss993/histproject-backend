package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// SetWebhook регистрирует вебхук бота в Telegram API.
// Если token или webhookURL пустые — вызов пропускается.
func SetWebhook(token, webhookURL string) {
	if token == "" || webhookURL == "" {
		return
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)

	params := url.Values{}
	params.Set("url", webhookURL)

	resp, err := (&http.Client{Timeout: 10 * time.Second}).
		PostForm(apiURL, params)
	if err != nil {
		log.Printf("[telegram] предупреждение: не удалось установить вебхук: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &result); err != nil || !result.OK {
		log.Printf("[telegram] предупреждение: setWebhook вернул ошибку: %s", result.Description)
		return
	}

	log.Printf("[telegram] вебхук успешно установлен: %s", webhookURL)
}
