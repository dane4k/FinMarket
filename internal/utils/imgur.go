package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func uploadImageToImgur(imageBytes []byte, accessToken string) (string, error) {
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	payload := map[string]string{"image": imageBase64}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %v, body: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	data := result["data"].(map[string]interface{})
	imageURL := data["link"].(string)

	return imageURL, nil
}
