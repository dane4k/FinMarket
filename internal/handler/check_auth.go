package handler

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckStatusHandler(c *gin.Context) {
	token := c.Param("token")
	status, cookie, err := service.CheckAuthStatus(token)
	if err != nil {
		if errors.Is(err, default_error.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "auth token not found"})
		} else if errors.Is(err, default_error.ErrTokenExpired) {
			c.JSON(http.StatusGone, gin.H{"error": "auth token expired"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		}
		return
	}

	if cookie != nil {
		c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}
