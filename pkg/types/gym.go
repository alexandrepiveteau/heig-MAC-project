package types

import (
	"context"
	"errors"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
	neo4jDriver neo4j.Driver,
) (primitive.ObjectID, error) {

	// 1. Create in MongoDB
	id, err := g.createInMongo(db)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	// 2. Create in Neo4j
	err = g.createInNeo4j(neo4jDriver)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	// Return mongo's id
	return id, nil
}

func (g *Gym) createInMongo(
	db *mongo.Database,
) (primitive.ObjectID, error) {

	id, err := gymCollection(db).InsertOne(
		context.TODO(),
		bson.D{
			{Key: "name", Value: g.Name},
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

func (g *Gym) createInNeo4j(driver neo4j.Driver) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := "CREATE (g:Gym) SET g.name = $name RETURN g"
		params := map[string]interface{}{
			"name": g.Name,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
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
