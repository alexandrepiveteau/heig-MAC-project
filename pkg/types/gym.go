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
) (string, error) {

	// 1. Create in MongoDB
	id, err := g.createInMongo(db)
	if err != nil {
		return "", err
	}

	// 2. Create in Neo4j
	err = g.createInNeo4j(neo4jDriver, id)
	if err != nil {
		return "", err
	}

	// Return mongo's id
	return id, nil
}

func (g *Gym) createInMongo(
	db *mongo.Database,
) (string, error) {

	id, err := gymCollection(db).InsertOne(
		context.TODO(),
		bson.D{{Key: "name", Value: g.Name}},
	)
	if err != nil {
		return "", err
	}

	// Assert type ObjectID
	objectId, ok := id.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("ObjectID was not found.")
	}

	return objectId.Hex(), nil
}

func (g *Gym) createInNeo4j(
	driver neo4j.Driver,
	gymId string,
) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := "CREATE (g:Gym) SET g = {name: $name, id: $id} RETURN g"
		params := map[string]interface{}{
			"id":   gymId,
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
) (string, error) {

	// Filter all gyms by name
	var res bson.M
	filter := bson.D{{Key: "name", Value: name}}
	err := gymCollection(db).FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return "", err
	}

	log.Printf("%+v\n", res)

	// Assert ObjectID type on _id
	id := res["_id"]
	objectId, ok := id.(primitive.ObjectID)
	if !ok {
		return "", errors.New("ObjectID was not found.")
	}

	return objectId.Hex(), nil
}
