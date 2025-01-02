package service

import (
	"errors"
	"fmt"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func generateJWT(userId int64, authRecordID uint) (string, error) {
	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":    jti,
		"userID": userId,
		"authID": authRecordID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	signedJWT, err := token.SignedString(jwtSecret)
	if err != nil {
		logrus.WithError(err).Errorf("%s for userID: %v", default_error.ErrSigningJWT, userId)
		return "", err
	}

	if err := repository.PinJTI(authRecordID, jti); err != nil {
		return "", err
	}
	return signedJWT, nil
}

func ParseUIDFromJWT(signedJWT string) (int64, error) {
	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("%s: %v", default_error.ErrSigningMethod, token.Header["alg"])
			return nil, errors.New(default_error.ErrSigningMethod)
		}
		return jwtSecret, nil
	})

	if err != nil {
		logrus.WithError(err).Error(default_error.ErrParseJWT)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["userID"]; exists {
			if parsedID, ok := userID.(float64); ok {
				return int64(parsedID), nil
			}
			logrus.Errorf("%s: %v", default_error.ErrInvalidIDType, userID)
			return 0, fmt.Errorf(default_error.ErrInvalidIDType)
		}
		logrus.Errorf(default_error.ErrEmptyClaims)
		return 0, errors.New(default_error.ErrEmptyClaims)
	}

	logrus.Errorf("%s: %s", default_error.ErrorInvalidJWT, signedJWT)
	return 0, errors.New(default_error.ErrorInvalidJWT)
}

func ExtractJTI(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("%s: %v", default_error.ErrSigningMethod, token.Header["alg"])
			return nil, errors.New(default_error.ErrSigningMethod)
		}
		return jwtSecret, nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("%s: %s", default_error.ErrParsingJWT, tokenString)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if jti, exists := claims["jti"]; exists {
			return jti.(string), nil
		}
		logrus.Errorf("%s: %v", default_error.ErrEmptyClaims, claims)
		return "", errors.New(default_error.ErrEmptyClaims)
	}

	logrus.Warn(default_error.ErrInvalidClaims)
	return "", errors.New(default_error.ErrInvalidClaims)
}
