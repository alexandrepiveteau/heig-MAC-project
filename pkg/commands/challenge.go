package commands

import (
	"climb/pkg/commands/keyboards"
	"climb/pkg/types"
	"climb/pkg/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type challengeStage int

const (
	challengeInit challengeStage = iota
	challengeUsername
	challengeGym
	challengeRoute
	challengeEnd
)

type challengeState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	// Stage of the progress in the command
	stage challengeStage

	// internal data
	usernameChoices []keyboards.Choice
	username        *string
	gym             *string
	route           *string
}

func (s *challengeState) init(update tgbotapi.Update) {
	msg1 := tgbotapi.NewMessage(utils.GetChatId(&update), "Challenging someone you follow.")
	msg2 := tgbotapi.NewMessage(utils.GetChatId(&update), "Here are the people you follow (mocked). If you can't find them in this list, type in their @username. Who do you want to challenge?")

	// TODO: Get followed people from DB, remove "(mocked)" from msg2
	follow := []string{
		"@alexandrepiveteau",
		"@matt989253",
		"@glsubri",
	}

	for _, username := range follow {
		s.usernameChoices = append(s.usernameChoices, keyboards.Choice{Action: username, Label: username})
	}

	msg2.ReplyMarkup = keyboards.NewInlineKeyboard(s.usernameChoices, 1)

	s.bot.Send(msg1)
	s.bot.Send(msg2)
	s.stage = challengeUsername
}

func (s *challengeState) rcvUsername(update tgbotapi.Update) {
	data, present := utils.GetInlineKeyboardData(
		update,
		keyboards.GetActions(s.usernameChoices)...,
	)
	if !present {
		// Get username when not in custom keyboard
		username, present := utils.GetMessageData(update)
		if !present {
			return // ignore update
		}
		data = username
	}
	s.username = &data

	utils.RemoveInlineKeyboard(s.bot, &update)

	reply := fmt.Sprintf("In which gym is the route you want to challenge %s to?", *s.username)
	msg := tgbotapi.NewMessage(utils.GetChatId(&update), reply)

	s.bot.Send(msg)
	s.stage = challengeGym
}

func (s *challengeState) rcvGym(update tgbotapi.Update) {
	data, present := utils.GetMessageData(update)
	if !present {
		return // ignore update
	}

	s.gym = &data

	msg := tgbotapi.NewMessage(utils.GetChatId(&update), " What is the route's name?")

	s.bot.Send(msg)
	s.stage = challengeRoute
}

func (s *challengeState) rcvRoute(update tgbotapi.Update) bool {
	data, present := utils.GetMessageData(update)
	if !present {
		return false // ignore update
	}

	s.route = &data

	reply := fmt.Sprintf("Great! You challenged %s to climb %s in %s!", *s.username, *s.route, *s.gym)
	msg := tgbotapi.NewMessage(utils.GetChatId(&update), reply)

	s.bot.Send(msg)
	s.stage = challengeEnd
	return true
}

func ChallengeCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
) {

	state := challengeState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		stage: challengeInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			return
		case update := <-comm.Updates:
			switch state.stage {
			case challengeInit:
				state.init(update)
				break
			case challengeUsername:
				state.rcvUsername(update)
				break
			case challengeGym:
				state.rcvGym(update)
				break
			case challengeRoute:
				if result := state.rcvRoute(update); result {
					commandTermination <- struct{}{}
				}
				// TODO: Challenge other user
				break
			case challengeEnd:
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry I'm lost.")

				bot.Send(msg)
				break
			}
		}
	}
}
