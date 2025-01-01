package handler

import (
	"github.com/dane4k/FinMarket/internal/auth"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CheckStatusHandler(c *gin.Context) {
	token := c.Param("token")
	record, err := repository.GetAuthRecord(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "auth token not found"})
		return
	}

	now := time.Now().UTC().Add(3 * time.Hour)
	expires := record.ExpiresAt.Truncate(time.Second)
	if now.After(expires) {
		c.JSON(http.StatusGone, gin.H{"error": "auth token expired"})
		return
	}

	if record.Status == "confirmed" {
		jwtToken, err := auth.GenerateJWT(record.TgID, record.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}
		c.SetCookie("jwtToken", jwtToken, 3600*24, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"status": "confirmed"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
	}
}
