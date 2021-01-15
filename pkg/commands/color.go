package commands

import (
	"climb/pkg/comm"
	"fmt"
	"log"

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
	stage Stage
}

func (s *state) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What is your favourite color?")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = colorKBD

	s.bot.Send(msg)
	s.stage = FavouriteColor
}

func (s *state) favourite(update tgbotapi.Update) {
	log.Print(update.CallbackQuery.Data)

	chatId := update.CallbackQuery.Message.ChatID
	msgID := update.CallbackQuery.Message.MessageID
	text := fmt.Sprintf("Sweet, your favourite color is %s\n", s.favourite)

	//msg := tgbotapi.NewEditMessageText(chatId, msgID, text)
	/*
	 *  msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What is your favourite color?")
	 *  msg.ReplyToMessageID = update.Message.MessageID
	 *  msg.ReplyMarkup = colorKBD
	 *
	 *  s.bot.Send(msg)
	 *  s.stage = FavouriteColor
	 */
}

// Entrypoint of bot command
func Color(
	comm comm.Comm,
	bot *tgbotapi.BotAPI,
) {
	state := state{
		bot:   bot,
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
