package types

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// A Comm is a way to share what channels are used to communicate between two goroutines
type Comm struct {
	Updates     chan tgbotapi.Update
	StopCommand chan interface{}
}

// InitComm will initialize a Comm and return it
func InitComm() Comm {
	return Comm{
		Updates:     make(chan tgbotapi.Update),
		StopCommand: make(chan interface{}),
	}
}
