package types

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type UserData struct {
	Username string
	Channel  chan tgbotapi.Update
	ChatId   int64
}

func (u *UserData) RegisterInNeo4j(
	driver neo4j.Driver,
) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := "CREATE (u:User) SET u = {name: $name} RETURN u"
		params := map[string]interface{}{
			"name": u.Username,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

func (u *UserData) Follow(
	driver neo4j.Driver,
	followingUsername string,
) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (me:User) WHERE me.name = $username
							MATCH (them:User) WHERE them.name = $followingUsername
							CREATE (me)-[:FOLLOWS]->(them)
							RETURN me`

		params := map[string]interface{}{
			"username":          u.Username,
			"followingUsername": followingUsername,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

func (u *UserData) Unfollow(
	driver neo4j.Driver,
	unfollowingUsername string,
) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (me:User) WHERE me.name = $username
							MATCH (them:User) WHERE them.name = $unfollowingUsername
							MATCH (me)-[f:FOLLOWS]->(them)
							DELETE f`

		params := map[string]interface{}{
			"username":            u.Username,
			"unfollowingUsername": unfollowingUsername,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return transRes, nil
	})

	return err
}

func (u *UserData) GetFollowing(
	driver neo4j.Driver,
) ([]string, error) {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	names, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (me:User)-[:FOLLOWS]->(following) WHERE me.name = $username RETURN following`

		params := map[string]interface{}{
			"username": u.Username,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		result := make([]string, 0)
		collect, _ := transRes.Collect()
		for _, res := range collect {
			node := res.Values[0].(dbtype.Node)

			name, ok := node.Props["name"].(string)
			if !ok {
				return nil, errors.New("Could not cast name to string")
			}
			result = append(result, name)
		}

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	retValue, ok := names.([]string)
	if !ok {
		return nil, errors.New("Could not cast names to []string")
	}

	return retValue, nil
}

func (u *UserData) GetFollowerRecommendation(
	driver neo4j.Driver,
) ([]string, error) {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	names, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		cypher := `MATCH (me:User)-[:FOLLOWS]->()-[:FOLLOWS]->(following) WHERE me.name = $username RETURN following`

		params := map[string]interface{}{
			"username": u.Username,
		}

		transRes, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		result := make([]string, 0)
		collect, _ := transRes.Collect()
		for _, res := range collect {
			node := res.Values[0].(dbtype.Node)

			name, ok := node.Props["name"].(string)
			if !ok {
				return nil, errors.New("Could not cast name to string")
			}
			result = append(result, name)
		}

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	retValue, ok := names.([]string)
	if !ok {
		return nil, errors.New("Could not cast names to []string")
	}

	return retValue, nil
}
