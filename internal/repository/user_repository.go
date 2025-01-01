package repository

import (
	"errors"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/model"
	"github.com/dane4k/FinMarket/internal/util"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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

func GetUser(userID int) (model.User, error) {
	var existingUser model.User

	if err := db.DB.Where("tg_id = ?", userID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, nil
		}
		logrus.WithError(err).Errorf("%s: %s", default_error.GetUserContext, default_error.ErrDatabaseError)
		return model.User{}, err
	}
	return existingUser, nil
}

func updateUser(user model.User) error {
	if err := db.DB.Save(&user).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", default_error.UpdateUserContext, default_error.ErrDatabaseError)
		return err
	}
	return nil
}

func saveNewUser(user model.User) error {
	if err := db.DB.Create(&user).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", default_error.SaveUserContext, default_error.ErrDatabaseError)
		return err
	}
	return nil
}

func PutUser(bot *tgbotapi.BotAPI, user *tgbotapi.User) error {
	usr, err := GetUser(user.ID)
	if err != nil {
		return err
	}
	if usr.TgID == 0 {
		avatarURL := util.DownloadTGAvatar(bot, user.ID)
		newUser := model.User{
			TgID:       int64(user.ID),
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

		return updateUser(usr)
	}
	return nil

}

func UpdateAvatarPic(userID int) error {

	user, err := GetUser(userID)
	if err != nil || user.TgID == 0 {
		return err
	}

	avatarURL := util.DownloadTGAvatar(bot, userID)
	user.AvatarURI = avatarURL

	return updateUser(user)
}
