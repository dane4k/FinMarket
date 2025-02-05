package tg_bot

import (
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TelegramBot struct {
	bot      *tgbotapi.BotAPI
	authRepo pgdb.AuthRepository
	userRepo pgdb.UserRepository
}

func NewTGBot(cfg *config.Config, authRepo pgdb.AuthRepository, userRepo pgdb.UserRepository) (*TelegramBot, error) {
	botToken := cfg.Telegram.Token

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logrus.WithError(err).Fatal("Error starting telegram bot")
		return nil, err
	}
	logrus.Infof("Authorized tg_bot on @%s", bot.Self.UserName)

	return &TelegramBot{
		bot:      bot,
		authRepo: authRepo,
		userRepo: userRepo,
	}, nil
}

func (tb *TelegramBot) Start() {
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates := tb.bot.GetUpdatesChan(upd)

	for update := range updates {
		if update.Message != nil {
			tb.handleMessage(update.Message)
		}
	}
}
