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

// RemoveInlineKeyboard Edits an old message and removes the associated inline keyboard
func RemoveInlineKeyboard(
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) {

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		update.CallbackQuery.Message.Text,
	)

	bot.Send(msg)
}

// GetChatId Returns the id of the chat where the message appeared
func GetChatId(update *tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}

	return 0
}
