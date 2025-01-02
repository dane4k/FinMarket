package handler

import (
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthHandler(c *gin.Context) {
	authToken, authLink, err := service.GenerateAuthLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth link"})
		return
	}
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"authToken": authToken,
		"AuthURL":   authLink,
	})
}
