package main

import (
	"climb/pkg/comm"
	"climb/pkg/controller"
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const envDebug = "BOT_DEBUG"
const envToken = "BOT_TOKEN"
const envNeo4j = "BOT_NEO4J"
const envMongo = "BOT_MONGO"

func main() {
	debug := os.Getenv(envDebug)
	token := os.Getenv(envToken)
	neo4jHost := os.Getenv(envNeo4j)
	mongoHost := os.Getenv(envMongo)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Neo4J
	driver, err := neo4j.NewDriver(neo4jHost, neo4j.NoAuth())
	if err != nil {
		log.Panic(err)
	}
	defer driver.Close()

	// Mongo
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoHost))
	if err != nil {
		log.Panic(err)
	}

	// TODO : Properly handle context cancellation.
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(ctx)

	bot.Debug = debug == "true"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Start controller
	controller := controller.GetController(
		bot,
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Prepare a map of user -> chan update
	userForwarder := make(map[int]chan tgbotapi.Update)

	for update := range updates {

		var userId int
		if update.Message != nil {
			userId = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			userId = update.CallbackQuery.From.ID
		}

		channel, prs := userForwarder[userId]

		if !prs {
			newChannel := make(chan tgbotapi.Update)
			userForwarder[userId] = newChannel
			channel = newChannel

			go handleUser(
				controller,
				channel,
			)
		}

		channel <- update
	}
}

func handleUser(
	ctrl controller.Controller,
	updates <-chan tgbotapi.Update,
) {
	var forwarder *comm.Comm
	send := ctrl.GetSendChannel()

	for update := range updates {

		logReception(update)

		if update.Message != nil && update.Message.IsCommand() {

			// Clean up previous commands
			if forwarder != nil {
				forwarder.QuitCommand <- struct{}{}
				forwarder = nil
			}

			// Get new command started
			switch update.Message.Command() {
			case "color":
				comm := ctrl.InstantiateColorCmd()
				forwarder = &comm
				forwarder.Updates <- update
				break
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know about this...")
				msg.ReplyToMessageID = update.Message.MessageID

				send <- msg
			}
		} else if forwarder != nil {
			forwarder.Updates <- update
		}
	}
}

// logReception logs a message based on it's type. The update argument is the message to log.
func logReception(update tgbotapi.Update) {
	user := GetUser(&update)

	var text string
	if update.Message != nil {
		text = update.Message.Text
	} else if update.CallbackQuery != nil {
		text = update.CallbackQuery.Data
	}

	log.Printf("[%s] %s", user.String(), text)
}

// GetUser Returns the user who sent an update
func GetUser(update *tgbotapi.Update) *tgbotapi.User {
	var user *tgbotapi.User

	if update.Message != nil {
		user = update.Message.From
	} else if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
	}

	return user
}
