package commands

import (
	"climb/pkg/comm"
	"climb/pkg/commands/keyboards"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// State definition
// Makeshift enum: https://golang.org/ref/spec#Iota

type Stage int

const (
	Init Stage = iota
	FavouriteColor
	LeastFavouriteColor
	End
)

type state struct {
	bot *tgbotapi.BotAPI

	// Stage of the progress in the command
	stage Stage

	favouriteColor *string
	leastFavColor  *string
}

func (s *state) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What is your favourite color?")
	msg.ReplyMarkup = keyboards.Color

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
	text := "Ok, but what is your least favourite color?"

	msg := tgbotapi.NewEditMessageText(chatId, msgID, text)
	msg.ReplyMarkup = &keyboards.Color

	s.bot.Send(msg)

	s.stage = LeastFavouriteColor
}

func (s *state) leastFav(update tgbotapi.Update) {
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

	s.stage = End
}

// Entrypoint of bot command
func Color(
	comm comm.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
) {

	state := state{
		bot:            bot,
		stage:          Init,
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
			case Init:
				state.init(update)
				break
			case FavouriteColor:
				state.favourite(update)
				break
			case LeastFavouriteColor:
				state.leastFav(update)
				commandTermination <- struct{}{} // Inform that we have terminated
				break
			case End:
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I'm lost.")

				bot.Send(msg)
				break
			}
		}
	}
}
