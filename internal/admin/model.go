package admin

import "time"

// Admin — модель администратора в БД.
type Admin struct {
	ID           int       `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}

// LoginRequest — тело запроса на авторизацию.
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse — пара токенов (используется внутри сервиса).
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

// AdminSession — сессия администратора в БД.
type AdminSession struct {
	ID           int       `db:"id"`
	AdminID      int       `db:"admin_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}
