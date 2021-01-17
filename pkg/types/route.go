package types

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
	neo4jDriver neo4j.Driver,
) (primitive.ObjectID, error) {
	// 1. Store in mongodb
	id, err := r.createInMongo(db, neo4jDriver)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	// 2. Create in Neo4j
	err = r.createInNeo4j(neo4jDriver)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	// Return mongo's id
	return id, nil
}

func (r *Route) createInMongo(
	db *mongo.Database,
	neo4jDriver neo4j.Driver,
) (primitive.ObjectID, error) {

	// Get corresponding gym or create it
	gymId, err := GymGetId(db, r.Name)
	if err != nil {
		gym := Gym{
			Name: r.Gym,
		}
		gymId, err = gym.Store(db, neo4jDriver)
		if err != nil {
			return primitive.ObjectID{}, err
		}
	}

	// Add route
	id, err := db.Collection("routes").InsertOne(
		context.TODO(),
		bson.D{
			{Key: "gym", Value: gymId},
			{Key: "name", Value: r.Name},
			{Key: "grade", Value: r.Grade},
			{Key: "holds", Value: r.Holds},
			{Key: "setDate", Value: r.SetDate},
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

func (r *Route) createInNeo4j(driver neo4j.Driver) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `CREATE (r:Route)
							SET r = {
								name: $name,
								grade: $grade,
								holds: $holds
								}
							RETURN r`

		params := map[string]interface{}{
			"name":  r.Name,
			"grade": r.Grade,
			"holds": r.Holds,
			//"setDate": r.SetDate, TODO: add date
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}
