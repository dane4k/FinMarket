package handler

import (
	"github.com/dane4k/FinMarket/internal/auth"
	"github.com/dane4k/FinMarket/internal/model"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const logoutSub = "you are already logged out"

func LoadHomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "home page",
	})
}

func LogoutHandler(c *gin.Context) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
		return
	}

	jti, err := auth.ExtractJTI(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
	}

	err = repository.InvalidateAuthRecord(jti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
	}

	c.SetCookie("jwtToken", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func ShowProfileHandler(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	usr := user.(model.User)

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"userID":    usr.TgID,
		"username":  usr.TgUsername,
		"name":      usr.Name,
		"rating":    usr.Rating,
		"avatarPic": usr.AvatarURI,
		"regDate":   usr.RegDate,
	})
}

func UpdateAvatarHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = repository.UpdateAvatarPic(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Аватарка обновлена"})
}