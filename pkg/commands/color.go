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
	FavouriteColor
	LeastFavouriteColor
)

type state struct {
	bot *tgbotapi.BotAPI

	// Stage of the progress in the command
	stage          Stage
	favouriteColor *string
}

func (s *state) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What is your favourite color?")
	msg.ReplyMarkup = colorKBD

	s.bot.Send(msg)
	s.stage = FavouriteColor
}

func (s *state) favourite(update tgbotapi.Update) {
	// Update state with new information
	data := update.CallbackQuery.Data
	s.favouriteColor = &data

	// Prepare next message
	chatId := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID
	text := "Ok but what is your least favourite color?"

	msg := tgbotapi.NewEditMessageText(chatId, msgID, text)
	msg.ReplyMarkup = &colorKBD

	s.bot.Send(msg)

	s.stage = LeastFavouriteColor
}

// Entrypoint of bot command
func Color(
	comm comm.Comm,
	bot *tgbotapi.BotAPI,
) {
	state := state{
		bot:            bot,
		stage:          Init,
		favouriteColor: nil,
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
			case FavouriteColor:
				state.favourite(update)
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I'm lost.")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
				break
			}
		}
	}
}
