package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
