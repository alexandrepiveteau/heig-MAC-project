package commands

import (
	"climb/pkg/commands/keyboards"
	"climb/pkg/types"
	"climb/pkg/utils"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type addRouteStage int

const (
	addRouteInit addRouteStage = iota
	addRouteGym
	addRouteName
	addRouteGrade
	addRouteHolds
	addRouteEnd
)

type addRouteState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	// Stage of the progress in the command
	stage addRouteStage

	// internal data
	gym   *string
	name  *string
	grade *string
	holds *string
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
	msg.ReplyMarkup = keyboards.NewInlineKeyboard(keyboards.GradeChoices, 3)

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

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Thanks! We've added this route.")
	msg.ParseMode = tgbotapi.ModeMarkdown

	s.bot.Send(msg)
	s.stage = addRouteEnd
}

func (s *addRouteState) save() {
	route := types.Route{
		Gym:   *s.gym,
		Name:  *s.name,
		Grade: *s.grade,
		Holds: *s.holds,
	}

	log.Println("Saving route")

	_, err := route.Store(s.mongodb, s.neo4jDriver)
	if err != nil {
		log.Printf("Error saving Route: %s\n", err.Error())
	}
}

func AddRouteCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
) {

	state := addRouteState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		stage: addRouteInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			// Save data in db, then quit
			state.save()
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
