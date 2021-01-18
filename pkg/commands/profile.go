package commands

import (
	"climb/pkg/types"
	"climb/pkg/utils"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type profileStage int

const (
	profileInit profileStage = iota
	profileUsername
	profileEnd
)

type profileState struct {
	bot         *tgbotapi.BotAPI
	mongodb     *mongo.Database
	neo4jDriver neo4j.Driver

	// Stage of the progress in the command
	stage profileStage

	user         types.UserData
	currentUsers map[string]types.UserData

	// internal data
	username *string
}

func (s *profileState) init(update tgbotapi.Update) {
	msg1 := tgbotapi.NewMessage(utils.GetChatId(&update), "Searching user profile.")
	msg2 := tgbotapi.NewMessage(utils.GetChatId(&update), "What is the @username of the person you want to check out? Make sure he already contacted me at least once.")

	s.bot.Send(msg1)
	s.bot.Send(msg2)
	s.stage = profileUsername
}

func (s *profileState) rcvUsername(update tgbotapi.Update) bool {
	data, present := utils.GetMessageData(update)
	if !present {
		return false
	}

	user, prs := s.currentUsers[data]
	if !prs {
		msg := tgbotapi.NewMessage(
			utils.GetChatId(&update),
			fmt.Sprintf("@%s never contacted me. Send him this link to get him started: t.me/climbot", data),
		)
		s.bot.Send(msg)
		return false
	}
	s.username = &data

	followers, following, attempts, err := user.GetProfile(s.neo4jDriver)
	if err != nil {
		log.Println(err.Error())
	}

	msg := tgbotapi.NewMessage(
		utils.GetChatId(&update),
		fmt.Sprintf("@%s has %d followers and is following %d other people.\nThey have made %d attempts.", data, followers, following, attempts),
	)
	s.bot.Send(msg)

	s.stage = profileEnd
	return true
}

func ProfileCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	mongodb *mongo.Database,
	neo4jDriver neo4j.Driver,
	user types.UserData,
	currentUsers map[string]types.UserData,
) {
	state := profileState{
		bot:         bot,
		mongodb:     mongodb,
		neo4jDriver: neo4jDriver,

		user:         user,
		currentUsers: currentUsers,

		stage: profileInit,
	}

	for {
		select {
		case <-comm.StopCommand:
			return
		case update := <-comm.Updates:
			switch state.stage {
			case profileInit:
				state.init(update)
				break
			case profileUsername:
				if finish := state.rcvUsername(update); !finish {
					return
				}
				state.stage = profileEnd
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
