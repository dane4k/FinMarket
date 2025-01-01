package bot

import (
	"github.com/dane4k/FinMarket/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func StartTelegramBot() {
	botToken := os.Getenv("TG_BOT_TOKEN")

	if botToken == "" {
		logrus.Fatal("TG_BOT_TOKEN is not set in .env")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logrus.WithError(err).Fatal("Error starting Telegram Bot")
	}

	logrus.Infof("Authorized on bot account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := message.Text
	if strings.HasPrefix(text, "/start ") {
		token := strings.TrimPrefix(text, "/start ")
		processToken(bot, message.Chat.ID, token, message.From)
	} else if text == "Выйти на всех устройствах" {
		handleExit(bot, message.Chat.ID, message.From.ID)
	}
}
func processToken(bot *tgbotapi.BotAPI, chatID int64, token string, user *tgbotapi.User) {
	if recordStatus, err := repository.GetTokenStatus(token); err != nil || recordStatus != "pending" {
		sendMessage(bot, chatID, "Ссылка недействительна или уже использована.", false)
		return
	}

	if err := repository.ConfirmToken(token, user.ID); err != nil {
		sendMessage(bot, chatID, "Ошибка при подтверждении токена.", false)
		return
	}

	err := repository.PutUser(bot, user)
	if err != nil {
		sendMessage(bot, chatID, "Ошибка при сохранении пользователя в базу данных.", false)
	}

	sendMessage(bot, chatID, "Вы успешно авторизовались в FinMarket. Если Вы не знаете, зачем, нажмите на кнопку", true)
	return
}
func createExitKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Выйти на всех устройствах"},
	}
	return tgbotapi.NewReplyKeyboard(buttons)
}

func handleExit(bot *tgbotapi.BotAPI, chatID int64, userID int) {
	if err := repository.InvalidateAllTokens(int64(userID)); err != nil {
		logrus.WithError(err).Error("Error invalidating user sessions")
	}
	sendMessage(bot, chatID, "Вы вышли из системы на всех устройствах", false)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, msg string, addKeyboard bool) {
	message := tgbotapi.NewMessage(chatID, msg)
	if addKeyboard {
		keyboard := createExitKeyboard()
		message.ReplyMarkup = keyboard
	}
	if _, err := bot.Send(message); err != nil {
		logrus.WithError(err).Error("Error sending message")
		return
	}
}
