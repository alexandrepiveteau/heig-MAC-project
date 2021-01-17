package commands

import (
	"climb/pkg/types"
	"climb/pkg/utils"
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
	climbRouteRating
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
	rating      *int
}

func (s *climbRouteState) init(update tgbotapi.Update) {
	msg1 := tgbotapi.NewMessage(utils.GetChatId(&update), "Adding a new attempt to an existing route.")
	msg2 := tgbotapi.NewMessage(utils.GetChatId(&update), "In which gym are you climbing?")

	s.bot.Send(msg1)
	s.bot.Send(msg2)
	s.stage = climbRouteGym
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
			// TODO : Save data in db, then quit
			return
		case update := <-comm.Updates:
			switch state.stage {
			case climbRouteInit:
				state.init(update)
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry I'm lost.")
				_, _ = bot.Send(msg)
				break
			}
		}
	}
}
