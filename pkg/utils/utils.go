package utils

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// logReception logs a message based on it's type. The update argument is the message to log.
func LogReception(update tgbotapi.Update) {
	user := GetUser(&update)

	var text string
	if update.Message != nil {
		text = update.Message.Text
	} else if update.CallbackQuery != nil {
		text = update.CallbackQuery.Data
	}

	log.Printf("[%s] %s", user.String(), text)
}

// GetUser Returns the user who sent an update
func GetUser(update *tgbotapi.Update) *tgbotapi.User {
	var user *tgbotapi.User

	if update.Message != nil {
		user = update.Message.From
	} else if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
	}

	return user
}
