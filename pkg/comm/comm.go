package comm

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Comm struct {
	Updates chan tgbotapi.Update
	Quit    chan interface{}
}
