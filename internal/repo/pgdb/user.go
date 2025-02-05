package pgdb

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/imgur"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	GetUser(userID int64) (*entity.User, error)
	UpdateUser(user *entity.User) error
	SaveNewUser(user *entity.User) error
	PutUser(bot *tgbotapi.BotAPI, user *tgbotapi.User) error
	UpdateAvatarPic(userID int64) error
}

type userRepository struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewUserRepository(db *gorm.DB, cfg *config.Config) UserRepository {
	return &userRepository{db: db, cfg: cfg}
}

var bot *tgbotapi.BotAPI

func InitTGBot(cfg *config.Config) {
	botToken := cfg.Telegram.Token
	bot, _ = tgbotapi.NewBotAPI(botToken)
}

func (ur *userRepository) GetUser(userID int64) (*entity.User, error) {
	var existingUser *entity.User

	if err := ur.db.Where("tg_id = ?", userID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logrus.Error(err.Error())
		return nil, err
	}
	return existingUser, nil
}

func (ur *userRepository) UpdateUser(user *entity.User) error {
	return ur.db.Save(&user).Error
}

func (ur *userRepository) SaveNewUser(user *entity.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepository) PutUser(bot *tgbotapi.BotAPI, user *tgbotapi.User) error {
	if bot == nil || user == nil {
		return errors.New("bot or/and user cannot be nil")
	}
	usr, err := ur.GetUser(user.ID)
	if err != nil {
		return err
	}
	if usr == nil {
		avatarURL := imgur.DownloadTGAvatar(ur.cfg, bot, user.ID)
		newUser := &entity.User{
			TgID:       user.ID,
			Name:       user.FirstName + " " + user.LastName,
			TgUsername: user.UserName,
			AvatarURI:  avatarURL,
			RegDate:    time.Now(),
		}
		return ur.SaveNewUser(newUser)
	}
	newName := user.FirstName + " " + user.LastName

	if usr.Name != newName || usr.TgUsername != user.UserName {
		usr.Name = user.FirstName + " " + user.LastName
		usr.TgUsername = user.UserName

		return ur.UpdateUser(usr)
	}
	return nil

}

func (ur *userRepository) UpdateAvatarPic(userID int64) error {
	user, err := ur.GetUser(userID)
	if err != nil || user == nil {
		return errors.New("failed to update avatar picture")
	}
	avatarURL := imgur.DownloadTGAvatar(ur.cfg, bot, userID)
	user.AvatarURI = avatarURL

	return ur.UpdateUser(user)
}
