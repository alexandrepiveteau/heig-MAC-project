package main

import (
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
	neo4jDriver, err := neo4j.NewDriver(neo4jHost, neo4j.NoAuth())
	if err != nil {
		log.Panic(err)
	}
	defer neo4jDriver.Close()

	// Mongo
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoHost))
	if err != nil {
		log.Panic(err)
	}

	// TODO : Properly handle context cancellation.
	ctx := context.TODO()
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	bot.Debug = debug == "true"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Start controller
	controller := controller.GetController(
		bot,
		neo4jDriver,
		mongoClient,
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err.Error())
	}

	for update := range updates {
		controller.GetAssociatedChan(update) <- update
	}
}
