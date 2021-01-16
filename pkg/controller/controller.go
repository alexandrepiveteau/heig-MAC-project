package controller

import (
	"climb/pkg/comm"
	"climb/pkg/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Controller interface {
	GetSendChannel() chan<- tgbotapi.Chattable

	InstantiateColorCmd(commandTermination chan interface{}) comm.Comm
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

func (c *controller) InstantiateColorCmd(commandTermination chan interface{}) comm.Comm {
	comm := comm.Comm{
		Updates:     make(chan tgbotapi.Update),
		StopCommand: make(chan interface{}),
	}

	go commands.Color(comm, commandTermination, c.bot)

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
