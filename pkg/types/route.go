package types

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Route struct {
	Gym     string
	Name    string
	Grade   string
	Holds   string
	SetDate string
}

// Store will store a route in MongoDB correctly
func (r *Route) Store(
	db *mongo.Database,
) (*mongo.InsertOneResult, error) {

	// TODO: get Gym id before inserting
	id, err := db.Collection("routes").InsertOne(context.TODO(), r)
	if err != nil {
		log.Println(err.Error())
	}

	return id, err
}
