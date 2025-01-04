package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
	"github.com/sirupsen/logrus"
)

func GenerateAuthLink() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		logrus.WithError(err).Error(service_errs.ErrGeneratingLink)
		return "", "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	err = pgdb.CreateAuthRecord(token)
	if err != nil {
		logrus.WithError(err).Error(service_errs.ErrGeneratingLink)
		return "", "", err
	}

	return token, fmt.Sprintf("https://t.me/finmarket_auth_bot?start=%s", token), nil
}
