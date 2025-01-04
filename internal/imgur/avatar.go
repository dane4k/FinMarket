package imgur

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func DownloadTGAvatar(bot *tgbotapi.BotAPI, userID int64) string {
	var AccessToken = os.Getenv("IMGUR_ACCESS_TOKEN")
	var defaultAvatar = os.Getenv("DEFAULT_AVATAR")

	userPhotos, err := bot.GetUserProfilePhotos(
		tgbotapi.UserProfilePhotosConfig{
			UserID: userID,
		},
	)
	if err != nil {
		logrus.WithError(err).Error(ErrDownloadingPic)
		return ""
	}
	if len(userPhotos.Photos) == 0 {
		return defaultAvatar
	}

	latestPhoto := userPhotos.Photos[0]
	userPhoto := latestPhoto[len(latestPhoto)-1]
	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: userPhoto.FileID,
	})
	if err != nil {
		logrus.WithError(err).Error(ErrDownloadingPic)
		return ""
	}

	resp, err := http.Get(file.Link(bot.Token))
	if err != nil {
		logrus.WithError(err).Error(ErrDownloadingPic)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logrus.WithError(err).Error(ErrDownloadingPic)
		}
	}(resp.Body)

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error(ErrDownloadingPic)
		return ""
	}
	imgurURL, err := uploadImageToImgur(imageBytes, AccessToken)
	if err != nil {
		return ""
	}

	return imgurURL
}
