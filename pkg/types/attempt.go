package types

import (
	"context"
	"errors"
	"fmt"

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

	// 1. Store in mongodb
	id, err := a.createInMongo(db, neo4jDriver)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Creating in mongodb: %w", err)
	}

	// 2. Create in Neo4j
	err = a.createInNeo4j(neo4jDriver, id)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Creating in Neo4j: %w", err)
	}

	// 3. Link with Route
	gymId, err := GymGetId(db, a.GymName)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Retrieving gym to link: %w", err)
	}

	routeId, err := RouteGetId(db, gymId, a.RouteName)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Retrieving route to link: %w", err)
	}

	err = a.linkWith(neo4jDriver, id, routeId)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Linking in Neo4j: %w", err)
	}

	// Return mongo's id
	return id, nil
}

func (a *Attempt) createInMongo(
	db *mongo.Database,
	neo4jDriver neo4j.Driver,
) (primitive.ObjectID, error) {

	gymId, err := GymGetId(db, a.GymName)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Retrieving gymId: %w", err)
	}

	routeId, err := RouteGetId(db, gymId, a.RouteName)
	if err != nil {
		return primitive.NewObjectID(), fmt.Errorf("Retrieving routeId: %w", err)
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
		return primitive.NewObjectID(), fmt.Errorf("Inserting route: %w", err)
	}

	// Assert type ObjectID
	objectId, ok := id.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NewObjectID(), errors.New("ObjectID was not found.")
	}

	return objectId, nil
}

func (a *Attempt) createInNeo4j(
	driver neo4j.Driver,
	id primitive.ObjectID,
) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `CREATE (a:Attempt)
							SET a = {
							  id: $id,
								proposedGrade: $proposedGrade,
								performance: $performance,
								}
							RETURN a`

		params := map[string]interface{}{
			"id":            id.String(),
			"proposedGrade": a.ProposedGrade,
			"performance":   a.Performance,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

// linkWith Links an Attempt with a Route in Neo4j with the "" label
func (a *Attempt) linkWith(
	driver neo4j.Driver,
	routeId primitive.ObjectID,
	attemptId primitive.ObjectID,
) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (a:Attempt) WHERE a.id = $aId
							MATCH (r:Route) WHERE r.id = $rId
							CREATE (a)-[:TRY_TO_CLIMB]->(r)
							RETURN r`

		params := map[string]interface{}{
			"aId": attemptId.String(),
			"rId": routeId.String(),
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}
