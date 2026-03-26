package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Danzhking/secure-audit/services/api/internal/middleware"
)

type AuthHandler struct {
	jwtSecret string
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{jwtSecret: jwtSecret}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := authenticate(req.Username, req.Password)
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
		return
	}

	token, err := middleware.GenerateToken(h.jwtSecret, req.Username, role, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сформировать токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": 86400,
	})
}

func authenticate(username, password string) string {
	// Демо-пользователи для диплома; в продакшене — БД или внешний IdP
	users := map[string]struct{ password, role string }{
		"admin":    {"admin", "admin"},
		"analyst":  {"analyst", "viewer"},
		"operator": {"operator", "viewer"},
	}

	u, ok := users[username]
	if !ok || u.password != password {
		return ""
	}
	return u.role
}
