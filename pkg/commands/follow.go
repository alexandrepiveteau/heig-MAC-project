package commands

import (
	"climb/pkg/commands/keyboards"
	"climb/pkg/types"
	"climb/pkg/utils"
	"fmt"
	"log"

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

	user         types.UserData
	currentUsers map[string]types.UserData

	// internal data
	usernameChoices []keyboards.Choice
	username        *string
}

func (s *followState) init(update tgbotapi.Update) {

	recommendation, err := s.user.GetFollowerRecommendation(s.neo4jDriver)
	if err != nil {
		log.Printf("When receiving follower recommendation : %s\n", err.Error())
	}

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the @username the person you want to follow? \n\nHere are a few people you might know: ")

	for _, username := range recommendation {
		s.usernameChoices = append(s.usernameChoices, keyboards.Choice{Action: username, Label: username})
	}

	if len(s.usernameChoices) > 0 {
		msg.ReplyMarkup = keyboards.NewInlineKeyboard(s.usernameChoices, 1)
	}

	s.bot.Send(msg)
	s.stage = followUsername
}

func (s *followState) rcvUsername(update tgbotapi.Update) bool {
	data, present := utils.GetInlineKeyboardData(
		update,
		keyboards.GetActions(s.usernameChoices)...,
	)
	if !present {
		// Get username when not in custom keyboard
		username, present := utils.GetMessageData(update)
		if !present {
			return false // ignore update
		}
		data = username
	}

	var text string
	_, prs := s.currentUsers[data]
	if !prs {
		text = "The requested user does not exist. Ask them to join!"
	} else {
		s.username = &data
		text = fmt.Sprintf("You're now following @%s !", data)
	}

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), text)
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
	user types.UserData,
	currentUsers map[string]types.UserData,
) {
	state := followState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		stage: followInit,

		user:         user,
		currentUsers: currentUsers,
	}

	for {
		select {
		case <-comm.StopCommand:
			if state.username != nil {
				user.Follow(neo4jDriver, *state.username)
			}
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
				break
			default:
				msg := tgbotapi.NewMessage(utils.GetChatId(&update), "Sorry I'm lost")
				bot.Send(msg)
				break
			}
		}
	}
}
