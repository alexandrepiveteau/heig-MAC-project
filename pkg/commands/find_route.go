package commands

import (
	"climb/pkg/commands/keyboards"
	"climb/pkg/types"
	"climb/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type findRouteStage int

const (
	findRouteInit findRouteStage = iota
	findRouteGym
	findRouteGrade
	findRouteHolds
	findRouteEnd
)

type findRouteState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	// Stage of the progress in the command
	stage findRouteStage

	// internal data
	gym   *string
	grade *string
	holds *string
}

func (s *findRouteState) init(update tgbotapi.Update) {
	msg1 := tgbotapi.NewMessage(utils.GetChatId(&update), "Searching for routes.")
	msg2 := tgbotapi.NewMessage(utils.GetChatId(&update), "In which gym do you want to find the route?")

	s.bot.Send(msg1)
	s.bot.Send(msg2)
	s.stage = findRouteGym
}

func (s *findRouteState) rcvGym(update tgbotapi.Update) {
	data, present := utils.GetMessageData(update)
	if !present {
		return // ignore update
	}

	s.gym = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the grade of the route?")
	msg.ReplyMarkup = keyboards.NewInlineKeyboard(keyboards.GradeChoices, 3)

	s.bot.Send(msg)
	s.stage = findRouteGrade
}

func (s *findRouteState) rcvGrade(update tgbotapi.Update) {
	data, present := utils.GetInlineKeyboardData(
		update,
		keyboards.GetActions(keyboards.GradeChoices)...,
	)
	if !present {
		return // ignore update
	}

	s.grade = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What color are the holds?")
	msg.ReplyMarkup = keyboards.NewInlineKeyboard(keyboards.ColorChoices, 3)

	s.bot.Send(msg)
	s.stage = findRouteHolds
}

func (s *findRouteState) rcvHolds(update tgbotapi.Update) bool {
	data, present := utils.GetInlineKeyboardData(
		update,
		keyboards.GetActions(keyboards.ColorChoices)...,
	)
	if !present {
		return false // ignore update
	}

	s.holds = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Thanks! We're looking for this route")
	msg.ParseMode = tgbotapi.ModeMarkdown

	s.bot.Send(msg)
	s.stage = findRouteEnd
	return true
}

func FindRouteCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
) {

	state := findRouteState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		stage: findRouteInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			return
		case update := <-comm.Updates:
			switch state.stage {
			case findRouteInit:
				state.init(update)
				break
			case findRouteGym:
				state.rcvGym(update)
				break
			case findRouteGrade:
				state.rcvGrade(update)
				break
			case findRouteHolds:
				if finish := state.rcvHolds(update); finish {
					commandTermination <- struct{}{}
				}
				// TODO: Return found routes to user
				break
			case findRouteEnd:
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry I'm lost.")

				bot.Send(msg)
				break
			}
		}
	}
}
