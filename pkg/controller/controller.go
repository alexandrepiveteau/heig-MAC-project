package controller

import (
	"climb/pkg/commands"
	"climb/pkg/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller interface {
	Bot() *tgbotapi.BotAPI

	InstantiateColorCmd(commandTermination chan interface{}) types.Comm
	InstantiateStartCmd(commandTermination chan interface{}) types.Comm

	AvailableCommands() []types.CommandDefinition
}

type controller struct {
	bot         *tgbotapi.BotAPI
	neo4jDriver *neo4j.Driver
	mongoClient *mongo.Client

	availableCommands []types.CommandDefinition
}

func GetController(
	bot *tgbotapi.BotAPI,
	neo4jDriver *neo4j.Driver,
	mongoClient *mongo.Client,
) Controller {

	controller := controller{
		bot:         bot,
		neo4jDriver: neo4jDriver,
		mongoClient: mongoClient,
	}

	startCmd := types.CommandDefinition{
		Command:       "start",
		Description:   "The start command shows available commands",
		Instantiation: controller.InstantiateStartCmd,
	}

	colorCmd := types.CommandDefinition{
		Command:       "color",
		Description:   "The color command will ask for your favourite color.",
		Instantiation: controller.InstantiateColorCmd,
	}

	controller.availableCommands = append(
		controller.availableCommands,
		startCmd,
		colorCmd,
	)

	return &controller
}

func (c *controller) AvailableCommands() []types.CommandDefinition {
	return c.availableCommands
}

func (c *controller) InstantiateColorCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.ColorCmd(comm, commandTermination, c.bot)

	return comm
}

func (c *controller) InstantiateStartCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.StartCmd(comm, commandTermination, c.bot)

	return comm
}

func (c *controller) Bot() *tgbotapi.BotAPI {
	return c.bot
}
