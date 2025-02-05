package service

import (
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

type UserService struct {
	cfg      *config.Config
	userRepo pgdb.UserRepository
	authRepo pgdb.AuthRepository
}

func NewUserService(userRepo pgdb.UserRepository, authRepo pgdb.AuthRepository, cfg *config.Config) *UserService {
	return &UserService{userRepo: userRepo, authRepo: authRepo, cfg: cfg}
}

func (us *UserService) LogoutUser(token string) error {
	jti, err := ExtractJTI(us.cfg, token)
	if err != nil {
		return ErrInvalidToken
	}

	if err = us.authRepo.InvalidateAuthRecord(jti); err != nil {
		return ErrLogoutFault
	}

	return nil
}

func (us *UserService) GetUserProfile(c *gin.Context) (map[string]interface{}, error) {
	user, ok := c.Get("user")
	if !ok {
		return nil, ErrUnauthorized
	}

	usr, ok := user.(*entity.User)
	if !ok {
		return nil, ErrUnauthorized
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

func (us *UserService) UpdateUserAvatar(c *gin.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		logrus.Info(err)
		return ErrInvalidUserData
	}

	if us.userRepo.UpdateAvatarPic(int64(userID)) != nil {
		return ErrUpdatingAvatar
	}

	return nil
}
