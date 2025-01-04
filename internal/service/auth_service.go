package service

import (
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
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
	record, err := pgdb.GetAuthRecord(token)
	if err != nil {
		return "", nil, service_errs.ErrTokenNotFound
	}

	now := time.Now().UTC().Add(3 * time.Hour)
	if now.After(record.ExpiresAt) {
		return "", nil, service_errs.ErrTokenExpired
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

func AuthorizeUser(c *gin.Context) (*entity.User, error) {
	userID, valid := IsAuthed(c)
	if !valid {
		return nil, service_errs.ErrUnauthorized
	}

	user, err := pgdb.GetUser(userID)
	if err != nil || user == nil {
		return nil, service_errs.ErrUserNotFound
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

	isValid, err := pgdb.IsJTIValid(jti)
	if err != nil || !isValid {
		return 0, false
	}

	userID, err := ParseUIDFromJWT(token)
	if err != nil {
		return 0, false
	}

	return userID, true
}
