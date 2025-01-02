package service

import (
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/model"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"time"
)

type Cookie struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

func CheckAuthStatus(token string) (string, *Cookie, error) {
	record, err := repository.GetAuthRecord(token)
	if err != nil {
		return "", nil, default_error.ErrTokenNotFound
	}

	now := time.Now().UTC().Add(3 * time.Hour)
	if now.After(record.ExpiresAt) {
		return "", nil, default_error.ErrTokenExpired
	}

	if record.Status == "confirmed" {
		jwtToken, err := generateJWT(record.TgID, record.ID)
		if err != nil {
			return "", nil, err
		}
		return "confirmed",
			&Cookie{
				Name:     "jwtToken",
				Value:    jwtToken,
				MaxAge:   3600 * 24,
				Path:     "/",
				Domain:   "localhost",
				Secure:   false,
				HttpOnly: true,
			}, nil
	} else {
		return "pending", nil, nil
	}
}

func AuthorizeUser(c *gin.Context) (model.User, error) {
	userID, valid := IsAuthed(c)
	if !valid {
		return model.User{}, default_error.ErrUnauthorized
	}

	user, err := repository.GetUser(int(userID))
	if err != nil || user.TgID == 0 {
		return model.User{}, default_error.ErrUserNotFound
	}
	return user, nil
}

func IsAuthed(c *gin.Context) (int64, bool) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		return 0, false
	}

	jti, err := ExtractJTI(token)
	if err != nil {
		return 0, false
	}

	isValid, err := repository.IsJTIValid(jti)
	if err != nil || !isValid {
		return 0, false
	}

	userID, err := ParseUIDFromJWT(token)
	if err != nil {
		return 0, false
	}

	return userID, true
}
