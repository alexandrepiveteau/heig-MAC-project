package types

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
) (primitive.ObjectID, error) {

	// Get corresponding gym or create it
	gymId, err := GymGetId(db, r.Name)
	if err != nil {
		gym := Gym{
			Name: r.Gym,
		}
		gymId, err = gym.Store(db)
	}

	// Add route
	id, err := db.Collection("routes").InsertOne(
		context.TODO(),
		bson.D{
			{"gym", gymId},
			{"name", r.Name},
			{"grade", r.Grade},
			{"holds", r.Holds},
			{"setDate", r.SetDate},
		},
	)

	if err != nil {
		return primitive.NewObjectID(), err
	}

	// Assert type ObjectID
	objectId, ok := id.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NewObjectID(), errors.New("ObjectID was not found.")
	}

	return objectId, nil
}
