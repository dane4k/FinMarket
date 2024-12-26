package handlers

import (
	"github.com/dane4k/FinMarket/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AuthUser(c *gin.Context) {
	authToken, authLink, err := services.GenerateAuthLink()
	if err != nil {
		logrus.WithError(err).Error("Error generating auth link")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth link"})
		return
	}
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"authToken": authToken,
		"AuthURL":   authLink,
	})
}
