package types

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type UserData struct {
	Username string
	Channel  chan tgbotapi.Update
	ChatId   int64
}
