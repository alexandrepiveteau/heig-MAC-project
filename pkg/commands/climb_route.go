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

type climbRouteStage int

const (
	climbRouteInit climbRouteStage = iota
	climbRouteGym
	climbRouteRoute
	climbRoutePerformance
	climbRouteGrade
	climbRouteEnd
)

type climbRouteState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	// Stage of the progress in the command
	stage climbRouteStage

	// internal data
	gym         *string
	route       *string
	performance *string
	grade       *string
}

func (s *climbRouteState) init(update tgbotapi.Update) {
	msg1 := tgbotapi.NewMessage(utils.GetChatId(&update), "Adding a new attempt to an existing route.")
	msg2 := tgbotapi.NewMessage(utils.GetChatId(&update), "In which gym are you climbing?")

	_, _ = s.bot.Send(msg1)
	_, _ = s.bot.Send(msg2)
	s.stage = climbRouteGym
}

func (s *climbRouteState) rcvGym(update tgbotapi.Update) {
	data := update.Message.Text
	s.gym = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the name of the route?")

	_, _ = s.bot.Send(msg)
	s.stage = climbRouteRoute
}

func (s *climbRouteState) rcvRoute(update tgbotapi.Update) {
	data := update.Message.Text
	s.route = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What was your performance?")
	msg.ReplyMarkup = keyboards.Performance

	_, _ = s.bot.Send(msg)
	s.stage = climbRoutePerformance
}

func (s *climbRouteState) rcvPerformance(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	s.performance = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "How would you grade the route?")
	msg.ReplyMarkup = keyboards.Grade

	_, _ = s.bot.Send(msg)
	s.stage = climbRouteGrade
}

func (s *climbRouteState) rcvGrade(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	s.grade = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Long live the swollen forearms! 💪")
	_, _ = s.bot.Send(msg)
	s.stage = climbRouteEnd
}

func (s *climbRouteState) save() {
	attempt := types.Attempt{
		GymName:       *s.gym,
		RouteName:     *s.route,
		ProposedGrade: *s.grade,
		Performance:   *s.performance,
	}

	log.Printf("Saving attempt %+v\n", attempt)

	_, err := attempt.Store(s.mongodb, s.neo4jDriver)
	if err != nil {
		log.Printf("Error saving Attempt: %s\n", err.Error())
	}
}

func ClimbRouteCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
) {

	state := climbRouteState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		stage: climbRouteInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			// Save data in db, then quit
			state.save()
			return
		case update := <-comm.Updates:
			switch state.stage {
			case climbRouteInit:
				state.init(update)
				break
			case climbRouteGym:
				state.rcvGym(update)
				break
			case climbRouteRoute:
				state.rcvRoute(update)
				break
			case climbRoutePerformance:
				state.rcvPerformance(update)
				break
			case climbRouteGrade:
				state.rcvGrade(update)
				commandTermination <- struct{}{}
				break
			case climbRouteEnd:
				break
			default:
				id := utils.GetChatId(&update)
				msg := tgbotapi.NewMessage(id, "Sorry I'm lost.")
				_, _ = bot.Send(msg)
				break
			}
		}
	}
}
