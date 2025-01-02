package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/sirupsen/logrus"
)

func GenerateAuthLink() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		logrus.WithError(err).Error(default_error.ErrGeneratingLink)
		return "", "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	err = repository.CreateAuthRecord(token)
	if err != nil {
		logrus.WithError(err).Error(default_error.ErrGeneratingLink)
		return "", "", err
	}

	return token, fmt.Sprintf("https://t.me/finmarket_auth_bot?start=%s", token), nil
}
