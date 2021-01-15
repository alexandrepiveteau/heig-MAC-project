package commands

import (
	"climb/pkg/comm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Color(
	comm comm.Comm,
	send chan<- tgbotapi.Chattable,
) {
	for {
		select {
		case <-comm.Quit:
			return
		case update := <-comm.Updates:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "colors everywhere")
			msg.ReplyToMessageID = update.Message.MessageID

			send <- msg
		}
	}
}
