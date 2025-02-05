package handler

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

const logoutSub = "you are already logged out"

func (uh *UserHandler) LogoutHandler(c *gin.Context) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
		return
	}

	err = uh.userService.LogoutUser(token)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": logoutSub})
		}
	}

	c.SetCookie("jwtToken", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func (uh *UserHandler) ShowProfileHandler(c *gin.Context) {
	usrData, err := uh.userService.GetUserProfile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "profile.html", usrData)
}

func (uh *UserHandler) UpdateAvatarHandler(c *gin.Context) {
	err := uh.userService.UpdateUserAvatar(c)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserData) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Аватарка обновлена"})
}
