package handler

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) Authorize(c *gin.Context) {
	authToken, authLink, err := ah.authService.CreateAuthLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth link"})
		return
	}
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"authToken": authToken,
		"AuthURL":   authLink,
	})
}

func (ah *AuthHandler) CheckStatus(c *gin.Context) {
	token := c.Param("token")

	status, cookie, err := ah.authService.CheckAuthStatus(token)
	if err != nil {
		if errors.Is(err, service.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "auth token not found"})
			return
		} else if errors.Is(err, service.ErrTokenExpired) {
			c.JSON(http.StatusGone, gin.H{"error": "auth token expired"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
	}
	if cookie != nil {
		c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
