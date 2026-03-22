package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler — HTTP-обработчики аутентификации администраторов.
type Handler struct {
	svc *Service
}

// NewHandler создаёт handler администраторов.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Login авторизует администратора и устанавливает токены в куки.
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка валидации: " + err.Error()})
		return
	}

	tokens, err := h.svc.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	setTokenCookies(c, tokens)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Refresh читает refresh-токен из куки и выдаёт новую пару токенов.
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Отсутствует refresh-токен"})
		return
	}

	tokens, err := h.svc.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	setTokenCookies(c, tokens)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Logout удаляет сессию и очищает куки.
func (h *Handler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		_ = h.svc.Logout(c.Request.Context(), refreshToken)
	}

	clearTokenCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func setTokenCookies(c *gin.Context, tokens *TokenResponse) {
	c.SetCookie("access_token", tokens.AccessToken, 15*60, "/", "", false, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, 7*24*60*60, "/api/admin/", "", false, true)
}

func clearTokenCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/api/admin/", "", false, true)
}
