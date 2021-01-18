package controller

import (
	"climb/pkg/commands"
	"climb/pkg/types"
	"climb/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller interface {
	Bot() *tgbotapi.BotAPI
	MongoDB() *mongo.Database

	AvailableCommands() []types.CommandDefinition
	GetAssociatedChan(update tgbotapi.Update) chan tgbotapi.Update
}

type controller struct {
	bot         *tgbotapi.BotAPI
	neo4jDriver neo4j.Driver
	mongodb     *mongo.Database

	// contacted users
	users map[string]types.UserData

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

		users: make(map[string]types.UserData),
	}

	// Define allowed commands
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

	challengeCmd := types.CommandDefinition{
		Command:       "challenge",
		Description:   "The challenge command will allow you to challenge a user you follow to climb a route",
		Instantiation: controller.instantiateChallengeCmd,
	}

	addRouteCmd := types.CommandDefinition{
		Command:       "addRoute",
		Description:   "The addRoute command will allow you to create a new route",
		Instantiation: controller.instantiateAddRouteCmd,
	}

	climbRouteCmd := types.CommandDefinition{
		Command:       "climbRoute",
		Description:   "The climbRoute command will allow you to save an attempt",
		Instantiation: controller.instantiateClimbRouteCmd,
	}

	findRouteCmd := types.CommandDefinition{
		Command:       "findRoute",
		Description:   "The findRoute command will allow you to find the name of routes",
		Instantiation: controller.instantiateFindRouteCmd,
	}

	followCmd := types.CommandDefinition{
		Command:       "follow",
		Description:   "The follow will allow you to follow another user",
		Instantiation: controller.instantiateFollowCmd,
	}

	unfollowCmd := types.CommandDefinition{
		Command:       "unfollow",
		Description:   "The unfollow will allow you to stop following another user",
		Instantiation: controller.instantiateUnfollowCmd,
	}

	profileCmd := types.CommandDefinition{
		Command:       "profile",
		Description:   "The profile will allow you to see infos about an user, like best route climbed and follower numbers",
		Instantiation: controller.instantiateProfileCmd,
	}

	// Update allowed commands in controller
	controller.availableCommands = append(
		controller.availableCommands,
		startCmd,
		colorCmd,
		challengeCmd,
		addRouteCmd,
		climbRouteCmd,
		findRouteCmd,
		followCmd,
		unfollowCmd,
		profileCmd,
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

func (c *controller) GetAssociatedChan(update tgbotapi.Update) chan tgbotapi.Update {
	username := utils.GetUser(&update).String()

	data, prs := c.users[username]
	if !prs {
		userdata := types.UserData{
			Username: username,
			Channel:  make(chan tgbotapi.Update),
			ChatId:   utils.GetChatId(&update),
		}

		c.users[username] = userdata
		data = userdata
	}

	return data.Channel
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

func (c *controller) instantiateChallengeCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.ChallengeCmd(
		comm,
		commandTermination,
		c.bot,
		c.mongodb,
		c.neo4jDriver,
	)

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

func (c *controller) instantiateClimbRouteCmd(commandTermination chan interface{}) types.Comm {
	comm := types.InitComm()

	go commands.ClimbRouteCmd(
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

func (c *controller) instantiateFollowCmd(
	commandTermination chan interface{},
) types.Comm {
	comm := types.InitComm()

	go commands.FollowCmd(
		comm,
		commandTermination,
		c.bot,
		c.mongodb,
		c.neo4jDriver,
	)

	return comm
}

func (c *controller) instantiateUnfollowCmd(
	commandTermination chan interface{},
) types.Comm {
	comm := types.InitComm()

	go commands.UnfollowCmd(
		comm,
		commandTermination,
		c.bot,
		c.mongodb,
		c.neo4jDriver,
	)

	return comm
}

func (c *controller) instantiateProfileCmd(
	commandTermination chan interface{},
) types.Comm {
	comm := types.InitComm()

	go commands.ProfileCmd(
		comm,
		commandTermination,
		c.bot,
		c.mongodb,
		c.neo4jDriver,
	)

	return comm
}
