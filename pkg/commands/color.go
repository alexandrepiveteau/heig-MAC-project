package commands

import (
	"climb/pkg/commands/keyboards"
	"climb/pkg/types"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// State definition
// Makeshift enum: https://golang.org/ref/spec#Iota

type colorStage int

const (
	colorInit colorStage = iota
	colorFavouriteColor
	colorLeastFavouriteColor
	colorEnd
)

type colorState struct {
	bot *tgbotapi.BotAPI

	// Stage of the progress in the command
	stage colorStage

	favouriteColor *string
	leastFavColor  *string
}

func (s *colorState) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What is your favourite color?")
	msg.ReplyMarkup = keyboards.NewInlineKeyboard(keyboards.ColorChoices, 3)

	s.bot.Send(msg)
	s.stage = colorFavouriteColor
}

func (s *colorState) favourite(update tgbotapi.Update) {
	// Update state with new information
	data := update.CallbackQuery.Data
	s.favouriteColor = &data

	// Prepare next message
	chatId := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID
	text := "Ok, but what is your least favourite color?"

	msg := tgbotapi.NewEditMessageText(chatId, msgID, text)
	msg.ReplyMarkup = &(keyboards.Color)

	s.bot.Send(msg)

	s.stage = colorLeastFavouriteColor
}

func (s *colorState) leastFav(update tgbotapi.Update) {
	// Update state with new information
	data := update.CallbackQuery.Data
	s.leastFavColor = &data

	// Prepare next message
	chatId := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID
	text := fmt.Sprintf(
		"I got it! Your favourite color is %s and your least favourite one is %s!",
		*s.favouriteColor,
		*s.leastFavColor,
	)

	msg := tgbotapi.NewEditMessageText(chatId, msgID, text)

	s.bot.Send(msg)

	s.stage = colorEnd
}

// Entrypoint of bot command
func ColorCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
) {

	state := colorState{
		bot:            bot,
		stage:          colorInit,
		favouriteColor: nil,
		leastFavColor:  nil,
	}

	for {
		select {
		case <-comm.StopCommand:
			// For now, simply quit. Later, we'll want to add all the information in the db
			return

		case update := <-comm.Updates:
			switch state.stage {
			case colorInit:
				state.init(update)
				break
			case colorFavouriteColor:
				state.favourite(update)
				break
			case colorLeastFavouriteColor:
				state.leastFav(update)
				commandTermination <- struct{}{} // Inform that we have terminated
				break
			case colorEnd:
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I'm lost.")

				bot.Send(msg)
				break
			}
		}
	}
}
