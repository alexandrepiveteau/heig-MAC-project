package commands

import (
	"climb/pkg/comm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var colorKBD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Red", "red"),
		tgbotapi.NewInlineKeyboardButtonData("Green", "green"),
		tgbotapi.NewInlineKeyboardButtonData("Blue", "blue"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Yellow", "yellow"),
		tgbotapi.NewInlineKeyboardButtonData("Orange", "orange"),
		tgbotapi.NewInlineKeyboardButtonData("Gray", "gray"),
	),
)

func Color(
	comm comm.Comm,
	send chan<- tgbotapi.Chattable,
) {
	for {
		select {
		case <-comm.Quit:
			// For now, simply quit. Later, we'll want to add all the information in the db
			return

		case update := <-comm.Updates:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "colors everywhere")
			msg.ReplyToMessageID = update.Message.MessageID

			send <- msg
		}
	}
}
