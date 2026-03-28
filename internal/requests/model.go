package requests

import "time"

// RequestListItem — краткая информация о заявке для списка.
type RequestListItem struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// RequestDetail — полная информация о заявке для карточки.
type RequestDetail struct {
	ID               int       `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	Description      string    `db:"description" json:"description"`
	EventDate        string    `db:"event_date" json:"eventDate"`
	EventTypeID      int       `db:"event_type_id" json:"eventTypeId"`
	Email            string    `db:"email" json:"email"`
	TelegramUsername string    `db:"telegram_username" json:"telegramUsername"`
	Status           string    `db:"status" json:"status"`
	AdminComment     *string   `db:"admin_comment" json:"adminComment"`
	SiteURL          string    `db:"site_url" json:"siteUrl"`
	ScreenshotURL    string    `db:"screenshot_url" json:"screenshotUrl"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}
