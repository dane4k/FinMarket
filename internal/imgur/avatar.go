package imgur

import (
	"github.com/dane4k/FinMarket/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func DownloadTGAvatar(cfg *config.Config, bot *tgbotapi.BotAPI, userID int64) string {
	defaultAvatar := cfg.Imgur.DefaultAvatar

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
	imgurURL, err := UploadImageToImgur(imageBytes, cfg.Imgur.AccessToken)
	if err != nil {
		return ""
	}

	return imgurURL
}
