package service

import (
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

func LogoutUser(token string) error {
	jti, err := ExtractJTI(token)
	if err != nil {
		return service_errs.ErrInvalidToken
	}

	if err = pgdb.InvalidateAuthRecord(jti); err != nil {
		return service_errs.ErrLogoutFault
	}

	return nil
}

func GetUserProfile(c *gin.Context) (map[string]interface{}, error) {
	user, ok := c.Get("user")
	if !ok {
		return nil, service_errs.ErrUnauthorized
	}

	usr, ok := user.(*entity.User)
	if !ok {
		return nil, service_errs.ErrInvalidUserData
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
		return service_errs.ErrInvalidUserData
	}

	if pgdb.UpdateAvatarPic(int64(userID)) != nil {
		return service_errs.ErrUpdatingAvatar
	}

	return nil
}
