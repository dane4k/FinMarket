package repository

import (
	"database/sql"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/models"
	"github.com/dane4k/FinMarket/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var bot *tgbotapi.BotAPI

func SaveAuthRecord(token string) error {
	authRecord := models.AuthRecord{
		Token:     token,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	return db.DB.Create(&authRecord).Error
}

func GetTokenStatus(token string) (string, error) {
	var record *models.AuthRecord
	if err := db.DB.Where("token = ?", token).First(&record).Error; err != nil {
		return "", err
	}
	return record.Status, nil
}

func ConfirmToken(token string, tgId int) error {
	var record *models.AuthRecord
	if err := db.DB.Where("token = ?", token).First(&record).Error; err != nil {
		return err
	}
	record.Status = "confirmed"
	record.TgID = int64(tgId)
	return db.DB.Save(&record).Error
}

func GetAuthRecord(token string) (*models.AuthRecord, error) {
	var record *models.AuthRecord
	if err := db.DB.Where("token = ?", token).First(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func IsTokenValid(tokenJTI string) bool {
	var invalidJWT *models.InvalidJWT
	if err := db.DB.Where("jwt_token = ?", tokenJTI).First(&invalidJWT).Error; err != nil {
		return true
	}
	return false
}

func InvalidateAllTokens(userId int64) error {
	var tokens []sql.NullString
	err := db.DB.Model(&models.AuthRecord{}).Where("tg_id = ?", userId).Pluck("jwt", &tokens).Error
	if err != nil {
		return err
	}

	var invalidTokens []models.InvalidJWT
	for _, token := range tokens {
		if token.Valid && token.String != "" {
			invalidTokens = append(invalidTokens, models.InvalidJWT{
				JWTToken: token.String,
			})
		}
	}

	if len(invalidTokens) > 0 {
		err = db.DB.Create(&invalidTokens).Error
		if err != nil {
			return err
		}
	}

	return db.DB.Model(&models.AuthRecord{}).Where("tg_id = ?", userId).Update("jwt", nil).Error
}

func InvalidateAuthRecord(jti string) error {
	invalidJWT := models.InvalidJWT{
		JWTToken: jti,
	}
	return db.DB.Create(&invalidJWT).Error
}

func PutUser(bot *tgbotapi.BotAPI, user *tgbotapi.User) error {
	var existingUser models.User
	err := db.DB.Where("tg_id = ?", user.ID).First(&existingUser).Error

	if err != nil {
		avatarURL := utils.DownloadTGAvatar(bot, user.ID)
		newUser := models.User{
			TgID:       int64(user.ID),
			Name:       user.FirstName + " " + user.LastName,
			TgUsername: user.UserName,
			AvatarPic:  avatarURL,
			RegDate:    time.Now(),
		}

		if err = db.DB.Create(&newUser).Error; err != nil {
			logrus.WithError(err).Error("Error adding user to the database")
			return err
		}
		logrus.Infof("User added: %s", newUser.TgUsername)

	} else {
		existingUser.Name = user.FirstName + " " + user.LastName
		existingUser.TgUsername = user.UserName

		if err = db.DB.Save(&existingUser).Error; err != nil {
			logrus.WithError(err).Error("Error updating existing user in the database")
			return err
		}
		logrus.Infof("User updated: %s", existingUser.TgUsername)
	}
	return nil
}

func InitTGBot() {
	botToken := os.Getenv("TG_BOT_TOKEN")

	if botToken == "" {
		logrus.Fatal("TG_BOT_TOKEN is not set in .env")
	}
	bot, _ = tgbotapi.NewBotAPI(botToken)
}

func UpdateAvatarPic(userID int) error {
	var existingUser models.User
	err := db.DB.Where("tg_id = ?", userID).First(&existingUser).Error
	if err != nil {
		logrus.WithError(err).Error("Error updating profile pic")
		return err
	}
	avatarURL := utils.DownloadTGAvatar(bot, userID)

	existingUser.AvatarPic = avatarURL

	if err = db.DB.Save(&existingUser).Error; err != nil {
		logrus.WithError(err).Error("Error updating existing user in the database")
		return err
	}
	logrus.Infof("user pic updated: %v", userID)
	return nil
}
