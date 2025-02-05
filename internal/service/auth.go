package service

import (
	"fmt"
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

type AuthService struct {
	cfg      *config.Config
	authRepo pgdb.AuthRepository
	userRepo pgdb.UserRepository
}

func NewAuthService(authRepo pgdb.AuthRepository, userRepo pgdb.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg, authRepo: authRepo, userRepo: userRepo}
}

func (as *AuthService) CheckAuthStatus(token string) (string, *Cookie, error) {
	record, err := as.authRepo.GetAuthRecord(token)
	if err != nil {
		return "", nil, ErrTokenNotFound
	}

	now := time.Now().UTC().Add(3 * time.Hour)
	if now.After(record.ExpiresAt) {
		return "", nil, ErrTokenExpired
	}

	if record.Status == "confirmed" {
		jwtToken, err := as.GenerateJWT(as.cfg, record.TgID, record.ID)
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

func (as *AuthService) AuthorizeUser(c *gin.Context) (*entity.User, error) {
	userID, valid := as.IsAuthed(c)
	if !valid {
		return nil, ErrUnauthorized
	}

	user, err := as.userRepo.GetUser(userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (as *AuthService) IsAuthed(c *gin.Context) (int64, bool) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		return 0, false
	}

	jti, err := ExtractJTI(as.cfg, token)
	if err != nil {
		return 0, false
	}

	isValid, err := as.authRepo.IsJTIValid(jti)
	if err != nil || !isValid {
		return 0, false
	}

	userID, err := ParseUIDFromJWT(as.cfg, token)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func (as *AuthService) GenerateJWT(cfg *config.Config, userId int64, authRecordID uint) (string, error) {
	jwtSecret := []byte(cfg.Auth.JWTSecret)

	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":    jti,
		"userID": userId,
		"authID": authRecordID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	signedJWT, err := token.SignedString(jwtSecret)
	if err != nil {
		logrus.WithError(err).Errorf("%s for userID: %v", ErrSigningJWT, userId)
		return "", err
	}

	if err := as.authRepo.PinJTI(authRecordID, jti); err != nil {
		return "", err
	}
	return signedJWT, nil
}

func (as *AuthService) CreateAuthLink() (string, string, error) {
	token, link, err := GenerateAuthLink()

	err = as.authRepo.CreateAuthRecord(token)
	if err != nil {
		logrus.WithError(err).Error(ErrGeneratingLink)
		return "", "", err
	}

	return token, link, nil
}

func ParseUIDFromJWT(cfg *config.Config, signedJWT string) (int64, error) {
	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("%s: %v", ErrSigningMethod.Error(), token.Header["alg"])
			return nil, ErrSigningMethod
		}
		return []byte(cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		logrus.WithError(err).Error(ErrParseJWT)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["userID"]; exists {
			if parsedID, ok := userID.(float64); ok {
				return int64(parsedID), nil
			}
			logrus.Errorf("%s: %v", ErrInvalidIDType, userID)
			return 0, fmt.Errorf(ErrInvalidIDType)
		}
		logrus.Error(ErrEmptyClaims.Error())
		return 0, ErrEmptyClaims
	}

	logrus.Errorf("%s: %s", ErrorInvalidJWT.Error(), signedJWT)
	return 0, ErrorInvalidJWT
}

func ExtractJTI(cfg *config.Config, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("%s: %v", ErrSigningMethod.Error(), token.Header["alg"])
			return nil, ErrSigningMethod
		}
		return []byte(cfg.Auth.JWTSecret), nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("%s: %s", ErrParsingJWT, tokenString)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if jti, exists := claims["jti"]; exists {
			return jti.(string), nil
		}
		logrus.Errorf("%s: %v", ErrEmptyClaims.Error(), claims)
		return "", ErrEmptyClaims
	}

	logrus.Warn(ErrInvalidClaims.Error())
	return "", ErrInvalidClaims
}
