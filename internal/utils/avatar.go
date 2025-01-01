package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func DownloadTGAvatar(bot *tgbotapi.BotAPI, userID int) string {
	var AccessToken = os.Getenv("IMGUR_ACCESS_TOKEN")
	var defaultAvatar = os.Getenv("DEFAULT_AVATAR")

	userPhotos, err := bot.GetUserProfilePhotos(
		tgbotapi.UserProfilePhotosConfig{
			UserID: userID,
		},
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get user %v profile pics", userID)
		return ""
	}
	logrus.Info(len(userPhotos.Photos))
	if len(userPhotos.Photos) == 0 {
		logrus.Info(0)
		return defaultAvatar
	}

	latestPhoto := userPhotos.Photos[0]
	userPhoto := latestPhoto[len(latestPhoto)-1]
	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: userPhoto.FileID,
	})
	if err != nil {
		logrus.WithError(err).Error("failed to get user profile pic")
		return ""
	}

	resp, err := http.Get(file.Link(bot.Token))
	if err != nil {
		logrus.WithError(err).Error("failed to download user profile pic")
		return ""
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read image bytes")
		return ""
	}
	imgurURL, err := uploadImageToImgur(imageBytes, AccessToken)
	if err != nil {
		logrus.WithError(err).Error("failed to upload image to imgur")
		return ""
	}

	logrus.Infof("uploaded users avatar to imgur: %s", imgurURL)
	return imgurURL
}
