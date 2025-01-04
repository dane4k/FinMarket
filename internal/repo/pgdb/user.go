package pgdb

import (
	"errors"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/imgur"
	"github.com/dane4k/FinMarket/internal/repo/repoerrs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"time"
)

var bot *tgbotapi.BotAPI

func InitTGBot() {
	botToken := os.Getenv("TG_BOT_TOKEN")
	if botToken == "" {
		logrus.Fatal("TG_BOT_TOKEN is not set in .env")
	}
	bot, _ = tgbotapi.NewBotAPI(botToken)
}

func GetUser(userID int64) (*entity.User, error) {
	var existingUser *entity.User

	if err := db.DB.Where("tg_id = ?", userID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("%s: %s", repoerrs.GetUserContext, repoerrs.ErrDatabaseError)
		return nil, err
	}
	return existingUser, nil
}

func UpdateUser(user *entity.User) error {
	if err := db.DB.Save(&user).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", repoerrs.UpdateUserContext, repoerrs.ErrDatabaseError)
		return err
	}
	return nil
}

func saveNewUser(user *entity.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", repoerrs.SaveUserContext, repoerrs.ErrDatabaseError)
		return err
	}
	return nil
}

func PutUser(bot *tgbotapi.BotAPI, user *tgbotapi.User) error {
	if bot == nil || user == nil {
		return errors.New("bot or/and user cannot be nil")
	}
	usr, err := GetUser(user.ID)
	if err != nil {
		return err
	}
	if usr == nil {
		avatarURL := imgur.DownloadTGAvatar(bot, user.ID)
		newUser := &entity.User{
			TgID:       user.ID,
			Name:       user.FirstName + " " + user.LastName,
			TgUsername: user.UserName,
			AvatarURI:  avatarURL,
			RegDate:    time.Now(),
		}
		return saveNewUser(newUser)
	}
	newName := user.FirstName + " " + user.LastName

	if usr.Name != newName || usr.TgUsername != user.UserName {
		usr.Name = user.FirstName + " " + user.LastName
		usr.TgUsername = user.UserName

		return UpdateUser(usr)
	}
	return nil

}

func UpdateAvatarPic(userID int64) error {
	user, err := GetUser(userID)
	if err != nil || user == nil {
		return errors.New("failed to update avatar picture")
	}
	avatarURL := imgur.DownloadTGAvatar(bot, userID)
	user.AvatarURI = avatarURL

	return UpdateUser(user)
}
