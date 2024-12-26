package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/dane4k/FinMarket/internal/repository"
)

func GenerateAuthLink() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	err = repository.SaveAuthRecord(token)
	if err != nil {
		return "", "", err
	}

	return token, fmt.Sprintf("https://t.me/finmarket_auth_bot?start=%s", token), nil
}
