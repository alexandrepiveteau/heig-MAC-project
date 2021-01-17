package types

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
