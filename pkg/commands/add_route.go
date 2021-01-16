package commands

import (
	"climb/pkg/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type addRouteStage int

const (
	addRouteInit addRouteStage = iota
	addRouteGym
	addRouteName
	addRouteGrade
	addRouteHolds
	addRouteDate
	// TODO: add image of Route
	addRouteEnd
)

type addRouteState struct {
	bot *tgbotapi.BotAPI

	// Stage of the progress in the command
	stage addRouteStage
}

func (s *addRouteState) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteGym
}

func (s *addRouteState) addGym(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteName
}

func (s *addRouteState) addName(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteGrade
}

func (s *addRouteState) addGrade(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteHolds
}

func (s *addRouteState) addHolds(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteDate
}

func (s *addRouteState) addDate(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteEnd
}

func AddRouteCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
) {

	state := addRouteState{
		bot:   bot,
		stage: addRouteInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			// For now, simply quit. Later, we'll want to add all the information in the db
			return

		case update := <-comm.Updates:
			switch state.stage {
			case addRouteInit:
				state.init(update)
				break
			case addRouteGym:
				state.addGym(update)
				break
			case addRouteName:
				state.addName(update)
				break
			case addRouteGrade:
				state.addGrade(update)
				break
			case addRouteHolds:
				state.addHolds(update)
				break
			case addRouteDate:
				state.addDate(update)
				break
			case addRouteEnd:
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I'm lost.")

				bot.Send(msg)
				break
			}
		}
	}
}
