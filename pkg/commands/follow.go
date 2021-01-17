package commands

import (
	"climb/pkg/types"
	"climb/pkg/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type followStage int

const (
	followInit followStage = iota
	followUsername
	followEnd
)

type followState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	// Stage of the progress in the command
	stage followStage

	// internal data
	username *string
}

func (s *followState) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the @username the person you want to follow? Make sure he already contacted me at least once.")

	s.bot.Send(msg)
	s.stage = followUsername
}

func (s *followState) rcvUsername(update tgbotapi.Update) bool {
	data, present := utils.GetMessageData(update)
	if !present {
		return false
	}

	// TODO : Check whether this user exists in the database before moving to the next state.
	msg := tgbotapi.NewMessage(utils.GetChatId(&update), fmt.Sprintf("You're now following @%s !", data))
	s.bot.Send(msg)

	s.stage = followEnd
	return true
}

func FollowCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
) {
	state := followState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		stage: followInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			// TODO : Actually follow the user.
			return
		case update := <-comm.Updates:
			switch state.stage {
			case followInit:
				state.init(update)
				break
			case followUsername:
				if finish := state.rcvUsername(update); !finish {
					return
				}
				state.stage = followEnd
				commandTermination <- struct{}{}
			default:
				msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Sorry I'm lost")
				bot.Send(msg)
				break
			}
		}
	}
}
