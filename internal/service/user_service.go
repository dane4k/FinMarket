package service

import (
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/model"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

func LogoutUser(token string) error {
	jti, err := ExtractJTI(token)
	if err != nil {
		return default_error.ErrInvalidToken
	}

	if err = repository.InvalidateAuthRecord(jti); err != nil {
		return default_error.ErrLogoutFault
	}

	return nil
}

func GetUserProfile(c *gin.Context) (map[string]interface{}, error) {
	user, ok := c.Get("user")
	if !ok {
		return nil, default_error.ErrUnauthorized
	}

	usr, ok := user.(model.User)
	if !ok {
		return nil, default_error.ErrInvalidUserData
	}

	return map[string]interface{}{
		"userID":    usr.TgID,
		"username":  usr.TgUsername,
		"name":      usr.Name,
		"rating":    usr.Rating,
		"avatarPic": usr.AvatarURI,
		"regDate":   usr.RegDate,
	}, nil
}

func UpdateUserAvatar(c *gin.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		logrus.Info(err)
		return default_error.ErrInvalidUserData
	}

	if repository.UpdateAvatarPic(userID) != nil {
		return default_error.ErrUpdatingAvatar
	}

	return nil
}
