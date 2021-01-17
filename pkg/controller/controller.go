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
	MongoDB() *mongo.Database

	AvailableCommands() []types.CommandDefinition
}

type controller struct {
	bot         *tgbotapi.BotAPI
	neo4jDriver neo4j.Driver
	mongodb     *mongo.Database

	availableCommands []types.CommandDefinition
}

// GetController will return a Controller
func GetController(
	bot *tgbotapi.BotAPI,
	neo4jDriver neo4j.Driver,
	mongoClient *mongo.Client,
) Controller {

	// Setup controller
	controller := controller{
		bot:         bot,
		neo4jDriver: neo4jDriver,
		mongodb:     mongoClient.Database("db"),
	}

	// Define allowd commands
	startCmd := types.CommandDefinition{
		Command:       "start",
		Description:   "The start command shows available commands",
		Instantiation: controller.instantiateStartCmd,
	}

	colorCmd := types.CommandDefinition{
		Command:       "color",
		Description:   "The color command will ask for your favourite color.",
		Instantiation: controller.instantiateColorCmd,
	}

	addRouteCmd := types.CommandDefinition{
		Command:       "addRoute",
		Description:   "The addRoute will allow you to create a new route",
		Instantiation: controller.instantiateAddRouteCmd,
	}

	findRouteCmd := types.CommandDefinition{
		Command:       "findRoute",
		Description:   "The findRoute will allow you to find the name of routes",
		Instantiation: controller.instantiateFindRouteCmd,
	}

	// Update allowed commands in controller
	controller.availableCommands = append(
		controller.availableCommands,
		startCmd,
		colorCmd,
		addRouteCmd,
		findRouteCmd,
	)

	return &controller
}

// Controller functions

func (c *controller) Bot() *tgbotapi.BotAPI {
	return c.bot
}

func (c *controller) MongoDB() *mongo.Database {
	return c.mongodb
}

func (c *controller) AvailableCommands() []types.CommandDefinition {
	return c.availableCommands
}

// Private functions

func (c *controller) instantiateStartCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.StartCmd(comm, commandTermination, c.bot, c.availableCommands)

	return comm
}

func (c *controller) instantiateColorCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.ColorCmd(comm, commandTermination, c.bot)

	return comm
}

func (c *controller) instantiateAddRouteCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.AddRouteCmd(
		comm,
		commandTermination,
		c.bot,
		c.mongodb,
		c.neo4jDriver,
	)

	return comm
}

func (c *controller) instantiateFindRouteCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.FindRouteCmd(
		comm,
		commandTermination,
		c.bot,
		c.mongodb,
		c.neo4jDriver,
	)

	return comm
}
