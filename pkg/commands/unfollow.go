package commands

import (
	"climb/pkg/types"
	"climb/pkg/utils"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type unfollowStage int

const (
	unfollowInit unfollowStage = iota
	unfollowUsername
	unfollowEnd
)

type unfollowState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	user         types.UserData
	currentUsers map[string]types.UserData

	// Stage of the progress in the command
	stage unfollowStage

	// internal data
	username *string
}

func (s *unfollowState) init(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the @username you want to unfollow ?")

	s.bot.Send(msg)
	s.stage = unfollowUsername
}

func (s *unfollowState) rcvUsername(update tgbotapi.Update) bool {
	data, present := utils.GetMessageData(update)
	if !present {
		return false
	}

	var text string
	_, prs := s.currentUsers[data]
	if !prs {
		text = "The requested user was not found. Please retype the username."
		msg := tgbotapi.NewMessage(utils.GetChatId(&update), text)
		s.bot.Send(msg)
		return false
	} else {
		s.username = &data
		text = fmt.Sprintf("You're not following @%s anymore !", data)
	}

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), text)
	s.bot.Send(msg)

	s.stage = unfollowEnd
	return true
}

func UnfollowCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
	user types.UserData,
	currentUsers map[string]types.UserData,
) {
	state := unfollowState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		user:         user,
		currentUsers: currentUsers,

		stage: unfollowInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			user.Unfollow(neo4jDriver, *state.username)
			return
		case update := <-comm.Updates:
			switch state.stage {
			case unfollowInit:
				state.init(update)
				break
			case unfollowUsername:
				if finish := state.rcvUsername(update); !finish {
					return
				}
				state.stage = unfollowEnd
				commandTermination <- struct{}{}
				break
			default:
				msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Sorry I'm lost")
				bot.Send(msg)
				break
			}
		}
	}
}
