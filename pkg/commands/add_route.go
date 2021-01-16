package commands

import (
	"climb/pkg/commands/keyboards"
	"climb/pkg/types"
	"climb/pkg/utils"

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

	// internal data
	gym     *string
	name    *string
	grade   *string
	holds   *string
	setDate *string
}

func (s *addRouteState) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "In which gym would you like to add the route?")

	s.bot.Send(msg)
	s.stage = addRouteGym
}

func (s *addRouteState) rcvGym(update tgbotapi.Update) {
	data := update.Message.Text
	s.gym = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the name of the route?")

	s.bot.Send(msg)
	s.stage = addRouteName
}

func (s *addRouteState) rcvName(update tgbotapi.Update) {
	data := update.Message.Text
	s.name = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the difficulty of the route?")
	msg.ReplyMarkup = keyboards.Grade

	s.bot.Send(msg)
	s.stage = addRouteGrade
}

func (s *addRouteState) rcvGrade(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	s.grade = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What colors are the holds?")
	msg.ReplyMarkup = keyboards.Color

	s.bot.Send(msg)
	s.stage = addRouteHolds
}

func (s *addRouteState) rcvHolds(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	s.holds = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "When was the route set? _(DD-MM-YYYY)_")
	msg.ParseMode = tgbotapi.ModeMarkdown

	s.bot.Send(msg)
	s.stage = addRouteDate
}

func (s *addRouteState) rcvDate(update tgbotapi.Update) {
	data := update.Message.Text
	s.setDate = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Thanks! We've added this route.")

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
				state.rcvGym(update)
				break
			case addRouteName:
				state.rcvName(update)
				break
			case addRouteGrade:
				state.rcvGrade(update)
				break
			case addRouteHolds:
				state.rcvHolds(update)
				break
			case addRouteDate:
				state.rcvDate(update)
				commandTermination <- struct{}{}
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
