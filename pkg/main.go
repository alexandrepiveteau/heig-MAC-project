package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson"
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Prepare a map of user -> chan update
	userForwarder := make(map[int]chan tgbotapi.Update)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		userId := update.Message.From.ID
		channel, prs := userForwarder[userId]

		if !prs {
			newChannel := make(chan tgbotapi.Update)
			userForwarder[userId] = newChannel
			channel = newChannel

			go handleUser(
				channel,
				bot,
				client,
				ctx,
				driver,
			)
		}

		channel <- update
	}
}

func handleUser(
	updates <-chan tgbotapi.Update,
	bot *tgbotapi.BotAPI,
	client *mongo.Client,
	ctx context.Context,
	driver neo4j.Driver,
) {
	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		database := client.Database("db")
		messages := database.Collection("messages")
		messages.InsertOne(ctx, bson.D{
			{Key: "body", Value: update.Message.Text},
		})
		count, err := messages.CountDocuments(ctx, bson.D{})

		reply := fmt.Sprintf("%d %s", count, update.Message.Text)

		if err != nil {
			reply = err.Error()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		// Create a placeholder session, to test Neo4j connectivity.
		session := driver.NewSession(neo4j.SessionConfig{})
		_ = session.Close()

		_, _ = bot.Send(msg)
	}
}
