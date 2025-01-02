package middleware

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := service.AuthorizeUser(c)
		if err != nil {
			if errors.Is(err, default_error.ErrUnauthorized) {
				c.Redirect(http.StatusFound, "/auth")
			} else if errors.Is(err, default_error.ErrUserNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				// нужны тесты на user not found
				//c.Redirect(http.StatusFound, "/profile")
			}
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
