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

// State definition
// Makeshift enum: https://golang.org/ref/spec#Iota

type Stage int

const (
	Init Stage = iota
)

type state struct {
	// Channel where to send messages
	send chan<- tgbotapi.Chattable

	// Stage of the progress in the command
	stage Stage
}

func (s *state) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Colors everywhere")
	msg.ReplyToMessageID = update.Message.MessageID

	s.send <- msg
}

// Entrypoint of bot command
func Color(
	comm comm.Comm,
	send chan<- tgbotapi.Chattable,
) {
	state := state{
		send:  send,
		stage: Init,
	}

	for {
		select {
		case <-comm.Quit:
			// For now, simply quit. Later, we'll want to add all the information in the db
			return

		case update := <-comm.Updates:
			switch state.stage {
			case Init:
				state.init(update)
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I'm lost.")
				msg.ReplyToMessageID = update.Message.MessageID

				send <- msg
				break
			}
		}
	}
}
