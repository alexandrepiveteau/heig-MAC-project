package types

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// gymCollection gives access to the gym collection in the database
func gymCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("gym")
}

type Gym struct {
	Name string
}

// Store adds the Gym object in mongodb and returns it's id.
//
// If the object already exists, it will still be added. All verifications
// should be made by the caller
func (g *Gym) Store(
	db *mongo.Database,
) (primitive.ObjectID, error) {

	id, err := gymCollection(db).InsertOne(
		context.TODO(),
		bson.D{
			{"name", g.Name},
		},
	)
	if err != nil {
		log.Println(err.Error())
	}

	// Assert type ObjectID
	objectId, ok := id.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NewObjectID(), errors.New("ObjectID was not found.")
	}

	return objectId, nil
}

// GymGetId returns the id of a Gym named name if it exists or an error
func GymGetId(
	db *mongo.Database,
	name string,
) (primitive.ObjectID, error) {

	// Filter all gyms by name
	filterCursor, err := gymCollection(db).Find(context.TODO(), bson.M{"name": "Le cube"})
	if err != nil {
		log.Fatal(err)
	}

	var gymsFiltered []bson.M
	if err = filterCursor.All(context.TODO(), &gymsFiltered); err != nil {
		log.Fatal(err)
	}

	if len(gymsFiltered) == 0 {
		return primitive.NewObjectID(), errors.New("Empty res")
	}

	// Cast result to ObjectID
	id := gymsFiltered[0]["_id"]

	objectId, ok := id.(primitive.ObjectID)
	if !ok {
		return primitive.NewObjectID(), errors.New("ObjectID was not found.")
	}

	return objectId, nil
}
