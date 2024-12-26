package auth

import (
	"errors"
	"fmt"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(userId int64, authRecordID uint) (string, error) {
	jti := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":    jti,
		"userID": userId,
		"authID": authRecordID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	signedJWT, err := token.SignedString(jwtSecret)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to sign JWT for userID:%v", userId)
		return "", err
	}

	err = db.DB.Model(&models.AuthRecord{}).Where("id = ?", authRecordID).Update("jwt", jti).Error
	if err != nil {
		logrus.WithError(err).Errorf("Failed to add jti to auth record: %v", authRecordID)
		return "", err
	}
	return signedJWT, nil
}

func ParseUIDFromJWT(signedJWT string) (int64, error) {
	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		logrus.WithError(err).Error("Failed to parse JWT")
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["userID"]; exists {
			if parsedID, ok := userID.(float64); ok {
				return int64(parsedID), nil
			}
			logrus.Errorf("userID is not of type float64: %v", userID)
			return 0, fmt.Errorf("userID is not of type float64")
		}
		logrus.Warn("userID not found in token claims")
		return 0, errors.New("userID not found in token claims")
	}

	logrus.Warn("Invalid token")
	return 0, errors.New("invalid token")
}

func ExtractJTI(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		logrus.WithError(err).Error("Failed to parse JWT")
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if jti, exists := claims["jti"]; exists {
			return jti.(string), nil
		}
		logrus.Warn("JTI not found in token claims")
		return "", errors.New("JTI not found in token claims")
	}

	logrus.Warn("Invalid token claims")
	return "", errors.New("invalid token claims")
}
