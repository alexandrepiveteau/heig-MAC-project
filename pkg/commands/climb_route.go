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
	data, present := utils.GetMessageData(update)
	if !present {
		return // ignore update
	}
	s.gym = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the name of the route?")

	_, _ = s.bot.Send(msg)
	s.stage = climbRouteRoute
}

func (s *climbRouteState) rcvRoute(update tgbotapi.Update) {
	data, present := utils.GetMessageData(update)
	if !present {
		return // ignore update
	}
	s.route = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What was your performance?")
	msg.ReplyMarkup = keyboards.NewInlineKeyboard(
		keyboards.PerformanceChoices,
		keyboards.SingleLine,
	)

	_, _ = s.bot.Send(msg)
	s.stage = climbRoutePerformance
}

func (s *climbRouteState) rcvPerformance(update tgbotapi.Update) {
	data, present := utils.GetInlineKeyboardData(
		update,
		keyboards.GetActions(keyboards.PerformanceChoices)...,
	)
	if !present {
		return // ignore update
	}
	s.performance = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "How would you grade the route?")
	msg.ReplyMarkup = keyboards.NewInlineKeyboard(keyboards.GradeChoices, 3)

	_, _ = s.bot.Send(msg)
	s.stage = climbRouteGrade
}

func (s *climbRouteState) rcvGrade(update tgbotapi.Update) bool {
	data, present := utils.GetInlineKeyboardData(
		update,
		keyboards.GetActions(keyboards.GradeChoices)...,
	)
	if !present {
		return false // ignore update
	}
	s.grade = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Long live the swollen forearms! 💪")
	_, _ = s.bot.Send(msg)
	s.stage = climbRouteEnd
	return true
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
				if finish := state.rcvGrade(update); finish {
					commandTermination <- struct{}{}
				}
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
