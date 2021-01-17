package types

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// attemptCollection gives access to the gym collection in the database
func attemptCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("attempts")
}

type Attempt struct {
	GymName       string
	RouteName     string
	ProposedGrade string
	Performance   string
}

func (a *Attempt) Store(
	db *mongo.Database,
	neo4jDriver neo4j.Driver,
) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), errors.New("Not implemented")
}

func (a *Attempt) createInMongo(
	db *mongo.Database,
	neo4jDriver neo4j.Driver,
) (primitive.ObjectID, error) {

	gymId, err := GymGetId(db, a.GymName)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	routeId, err := RouteGetId(db, gymId, a.RouteName)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	// Add route
	id, err := routeCollection(db).InsertOne(
		context.TODO(),
		bson.D{
			{Key: "gym", Value: gymId},
			{Key: "route", Value: routeId},
			{Key: "proposedGrade", Value: a.ProposedGrade},
			{Key: "performance", Value: a.Performance},
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
