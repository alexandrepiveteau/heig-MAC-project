package types

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// routeCollection gives access to the gym collection in the database
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
	creatingUser UserData,
) (string, error) {

	// 1. Store in mongodb
	routeId, err := r.createInMongo(db, neo4jDriver)
	if err != nil {
		return "", err
	}

	// 2. Create in Neo4j
	err = r.createInNeo4j(neo4jDriver, routeId)
	if err != nil {
		return "", err
	}

	// 3. Link with Gym
	gymId, err := GymGetId(db, r.Gym)
	if err != nil {
		return "", err
	}

	err = r.linkWith(neo4jDriver, routeId, gymId)
	if err != nil {
		return "", err
	}

	// 4. Link with Creating user
	err = r.linkWithCreatingUser(neo4jDriver, routeId, creatingUser.Username)
	if err != nil {
		return "", err
	}

	// Return mongo's id
	return routeId, nil
}

func (r *Route) createInMongo(
	db *mongo.Database,
	neo4jDriver neo4j.Driver,
) (string, error) {

	// Get corresponding gym or create it
	gymId, err := GymGetId(db, r.Gym)
	if err != nil {
		fmt.Printf("%s\n", err.Error())

		gym := Gym{
			Name: r.Gym,
		}
		gymId, err = gym.Store(db, neo4jDriver)
		if err != nil {
			return "", err
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
		return "", err
	}

	// Assert type ObjectID
	objectId, ok := id.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("ObjectID was not found.")
	}

	return objectId.Hex(), nil
}

func (r *Route) createInNeo4j(
	driver neo4j.Driver,
	routeId string,
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
			"id":    routeId,
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
	routeId string,
	gymId string,
) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (r:Route) WHERE r.id = $rId
							MATCH (g:Gym) WHERE g.id = $gId
							CREATE (r)-[:IS_IN]->(g)
							RETURN r`

		params := map[string]interface{}{
			"rId": routeId,
			"gId": gymId,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

func (r *Route) linkWithCreatingUser(
	driver neo4j.Driver,
	routeId string,
	username string,
) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (r:Route) WHERE r.id = $rId
							MATCH (u:User) WHERE u.name = $username
							CREATE (u)-[:CREATED]->(r)
							RETURN r`

		params := map[string]interface{}{
			"rId":      routeId,
			"username": username,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

// RouteGetId returns the id of a Route named name if it exists or an error
//
// gymId should be the id of the gym in which the route is
// name is the name of the route
func RouteGetId(
	db *mongo.Database,
	gymId string,
	name string,
) (string, error) {

	// Filter all routes by name and gymId
	filterCursor, err := routeCollection(db).Find(
		context.TODO(),
		bson.D{
			{Key: "gym", Value: gymId},
			{Key: "name", Value: name},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	var routesFiltered []bson.M
	if err = filterCursor.All(context.TODO(), &routesFiltered); err != nil {
		log.Fatal(err)
	}

	if len(routesFiltered) == 0 {
		return "", errors.New("Empty result, no route found")
	}

	// Cast result to ObjectID
	id := routesFiltered[0]["_id"]

	objectId, ok := id.(primitive.ObjectID)
	if !ok {
		return "", errors.New("ObjectID was not found.")
	}

	return objectId.Hex(), nil
}

// RouteFind retuns a slice of Route corresponding to the given parameters
func RouteFind(
	db *mongo.Database,
	gymName string,
	routeGrade string,
	routeHolds string,
) ([]Route, error) {

	gymId, err := GymGetId(db, gymName)
	if err != nil {
		return make([]Route, 0), fmt.Errorf("Getting gym id: %w", err)
	}

	// Filter all routes
	filterCursor, err := routeCollection(db).Find(
		context.TODO(),
		bson.D{
			{Key: "gym", Value: gymId},
			{Key: "grade", Value: routeGrade},
			{Key: "holds", Value: routeHolds},
		},
	)
	if err != nil {
		return make([]Route, 0), fmt.Errorf("Filtering collection: %w", err)
	}
	defer filterCursor.Close(context.TODO())

	routes := make([]Route, 0)
	for filterCursor.Next(context.TODO()) {
		var route Route
		if err = filterCursor.Decode(&route); err != nil {
			log.Println(err.Error())
			continue
		}

		routes = append(routes, route)
	}

	return routes, nil
}
