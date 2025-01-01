package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func uploadImageToImgur(imageBytes []byte, accessToken string) (string, error) {
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	payload := map[string]string{"image": imageBase64}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error(default_error.ErrUploadingPic)
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", bytes.NewReader(payloadBytes))
	if err != nil {
		logrus.WithError(err).Error(default_error.ErrUploadingPic)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logrus.WithError(err).Error(default_error.ErrUploadingPic)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logrus.WithError(err).Error(default_error.ErrUploadingPic)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error(default_error.ErrUploadingPic)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithError(err).Errorf("%s || status: %v, body: %s", default_error.ErrUploadingPic, resp.StatusCode, string(respBody))
		return "", fmt.Errorf(default_error.ErrUploadingPic)
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		logrus.WithError(err).Error(default_error.ErrUploadingPic)
		return "", err
	}

	data := result["data"].(map[string]interface{})
	imageURL := data["link"].(string)

	return imageURL, nil
}
