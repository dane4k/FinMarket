package middleware

import (
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService}
}

func (amw *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := amw.authService.AuthorizeUser(c)
		if err != nil {
			c.Set("isAuthed", false)
			c.Set("user", nil)
			c.Next()
			return
		}
		c.Set("isAuthed", true)
		c.Set("user", user)
		c.Next()
	}
}
