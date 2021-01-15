package controller

import (
	"climb/pkg/comm"
	"climb/pkg/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Controller interface {
	GetSendChannel() chan<- tgbotapi.Chattable

	InstantiateColorCmd() comm.Comm
}

type controller struct {
	bot      *tgbotapi.BotAPI
	sendChan chan tgbotapi.Chattable
}

func GetController(
	bot *tgbotapi.BotAPI,
) Controller {

	controller := controller{
		bot:      bot,
		sendChan: make(chan tgbotapi.Chattable),
	}

	go controller.startSender()

	return &controller
}

func (c *controller) InstantiateColorCmd() comm.Comm {
	comm := comm.Comm{
		Updates: make(chan tgbotapi.Update),
		Quit:    make(chan interface{}),
	}

	go commands.Color(comm, c.sendChan)

	return comm
}

func (c *controller) GetSendChannel() chan<- tgbotapi.Chattable {
	return c.sendChan
}

func (c *controller) startSender() {
	for msg := range c.sendChan {
		c.bot.Send(msg)
	}
}
