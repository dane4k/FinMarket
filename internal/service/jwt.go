package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
)

func GenerateAuthLink() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		logrus.WithError(err).Error(ErrGeneratingLink)
		return "", "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, fmt.Sprintf("https://t.me/finmarket_auth_bot?start=%s", token), nil
}
