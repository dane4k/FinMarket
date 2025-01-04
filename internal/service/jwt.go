package service

import (
	"fmt"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
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
		logrus.WithError(err).Errorf("%s for userID: %v", service_errs.ErrSigningJWT, userId)
		return "", err
	}

	if err := pgdb.PinJTI(authRecordID, jti); err != nil {
		return "", err
	}
	return signedJWT, nil
}

func ParseUIDFromJWT(signedJWT string) (int64, error) {
	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("%s: %v", service_errs.ErrSigningMethod.Error(), token.Header["alg"])
			return nil, service_errs.ErrSigningMethod
		}
		return jwtSecret, nil
	})

	if err != nil {
		logrus.WithError(err).Error(service_errs.ErrParseJWT)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["userID"]; exists {
			if parsedID, ok := userID.(float64); ok {
				return int64(parsedID), nil
			}
			logrus.Errorf("%s: %v", service_errs.ErrInvalidIDType, userID)
			return 0, fmt.Errorf(service_errs.ErrInvalidIDType)
		}
		logrus.Error(service_errs.ErrEmptyClaims.Error())
		return 0, service_errs.ErrEmptyClaims
	}

	logrus.Errorf("%s: %s", service_errs.ErrorInvalidJWT.Error(), signedJWT)
	return 0, service_errs.ErrorInvalidJWT
}

func ExtractJTI(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("%s: %v", service_errs.ErrSigningMethod.Error(), token.Header["alg"])
			return nil, service_errs.ErrSigningMethod
		}
		return jwtSecret, nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("%s: %s", service_errs.ErrParsingJWT, tokenString)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if jti, exists := claims["jti"]; exists {
			return jti.(string), nil
		}
		logrus.Errorf("%s: %v", service_errs.ErrEmptyClaims.Error(), claims)
		return "", service_errs.ErrEmptyClaims
	}

	logrus.Warn(service_errs.ErrInvalidClaims.Error())
	return "", service_errs.ErrInvalidClaims
}
