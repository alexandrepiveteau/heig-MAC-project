package comm

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// A Comm is a way to share what channels are used to communicate between two goroutines
type Comm struct {
	Updates     chan tgbotapi.Update
	StopCommand chan interface{}
}
