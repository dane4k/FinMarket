package handlers

import (
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/auth"
	"github.com/dane4k/FinMarket/internal/models"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "home page",
	})
}

func Logout(c *gin.Context) {

	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you are already logged out"})
		return
	}

	jti, err := auth.ExtractJTI(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you are already logged out"})
	}

	err = repository.InvalidateAuthRecord(jti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you are already logged out"})
	}

	c.SetCookie("jwtToken", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func ShowProfile(c *gin.Context) {
	token, err := c.Cookie("jwtToken")

	jti, err := auth.ExtractJTI(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/auth") // 402
	}

	if err != nil || !repository.IsTokenValid(jti) {
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	userID, err := auth.ParseUIDFromJWT(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	var user models.User
	if err := db.DB.Where("tg_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"userID":   user.TgID,
		"username": user.TgUsername,
		"name":     user.Name,
		"avatar":   user.AvatarURL,
		"rating":   user.Rating,
	})
}
