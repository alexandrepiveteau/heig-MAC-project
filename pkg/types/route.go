package types

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// gymCollection gives access to the gym collection in the database
func routeCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("routes")
}

type Route struct {
	Gym   string
	Name  string
	Grade string
	Holds string
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
	err = r.createInNeo4j(neo4jDriver, id)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	// 3. Link with Gym
	gymId, err := GymGetId(db, r.Gym)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	err = r.linkWith(neo4jDriver, id, gymId)

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
	id, err := routeCollection(db).InsertOne(
		context.TODO(),
		bson.D{
			{Key: "gym", Value: gymId},
			{Key: "name", Value: r.Name},
			{Key: "grade", Value: r.Grade},
			{Key: "holds", Value: r.Holds},
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

func (r *Route) createInNeo4j(
	driver neo4j.Driver,
	id primitive.ObjectID,
) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `CREATE (r:Route)
							SET r = {
							  id: $id,
								name: $name,
								grade: $grade,
								holds: $holds
								}
							RETURN r`

		params := map[string]interface{}{
			"id":    id.String(),
			"name":  r.Name,
			"grade": r.Grade,
			"holds": r.Holds,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

func (r *Route) linkWith(
	driver neo4j.Driver,
	routeId primitive.ObjectID,
	gymId primitive.ObjectID,
) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (r:Route) WHERE r.id = $rId
							MATCH (g:Gym) WHERE g.id = $gId
							CREATE (r)-[:IS_IN]->(g)
							RETURN r`

		params := map[string]interface{}{
			"rId": routeId.String(),
			"gId": gymId.String(),
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}