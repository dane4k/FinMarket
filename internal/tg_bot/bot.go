package tg_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"strings"
)

func (tb *TelegramBot) handleMessage(message *tgbotapi.Message) {
	text := message.Text
	if strings.HasPrefix(text, "/start ") {
		token := strings.TrimPrefix(text, "/start ")
		tb.processToken(message.Chat.ID, token, message.From)
	} else if text == "Выйти на всех устройствах" {
		tb.handleExit(message.Chat.ID, message.From.ID)
	}
}

func (tb *TelegramBot) processToken(chatID int64, token string, user *tgbotapi.User) {
	userID := user.ID
	record, err := tb.authRepo.GetAuthRecord(token)
	if err != nil || record.Status != "pending" {
		tb.sendMessage(chatID, "Ссылка недействительна или уже использована", false)
		return
	}

	if err = tb.authRepo.ConfirmToken(token, userID); err != nil {
		tb.sendMessage(chatID, "Ошибка при подтверждении токена.", false)
		return
	}

	err = tb.userRepo.PutUser(tb.bot, user)
	if err != nil {
		tb.sendMessage(chatID, "Ошибка при сохранении пользователя в базу данных.", false)
	}

	tb.sendMessage(chatID, "Вы успешно авторизовались в FinMarket. Если Вы не знаете, зачем, нажмите на кнопку", true)
	return
}

func (tb *TelegramBot) sendMessage(chatID int64, msg string, addKeyboard bool) {
	message := tgbotapi.NewMessage(chatID, msg)
	if addKeyboard {
		keyboard := tb.createExitKeyboard()
		message.ReplyMarkup = keyboard
	}
	if _, err := tb.bot.Send(message); err != nil {
		logrus.WithError(err).Error("Error sending message")
		return
	}
}

func (tb *TelegramBot) handleExit(chatID int64, userID int64) {
	if err := tb.authRepo.InvalidateAllTokens(userID); err != nil {
		tb.sendMessage(chatID, "У вас нет активных сессий", false)
	}
	tb.sendMessage(chatID, "Вы вышли из системы на всех устройствах", false)
}

func (tb *TelegramBot) createExitKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Выйти на всех устройствах"},
	}
	return tgbotapi.NewReplyKeyboard(buttons)
}
