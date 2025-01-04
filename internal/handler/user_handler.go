package handler

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

const logoutSub = "you are already logged out"

func LogoutHandler(c *gin.Context) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
		return
	}

	err = service.LogoutUser(token)
	if err != nil {
		if errors.Is(err, service_errs.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
		}
	}

	c.SetCookie("jwtToken", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func ShowProfileHandler(c *gin.Context) {
	usrData, err := service.GetUserProfile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "profile.html", usrData)
}

func UpdateAvatarHandler(c *gin.Context) {
	err := service.UpdateUserAvatar(c)
	if err != nil {
		if errors.Is(err, service_errs.ErrInvalidUserData) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Аватарка обновлена"})
}
