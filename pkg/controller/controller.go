package controller

import (
	"climb/pkg/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Controller interface {
	GetSendChannel() chan<- tgbotapi.Chattable

	InstantiateColorCmd() chan tgbotapi.Update
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

func (c *controller) InstantiateColorCmd() chan tgbotapi.Update {
	updates := make(chan tgbotapi.Update)

	go commands.Color(updates, c.sendChan)

	return updates
}

func (c *controller) GetSendChannel() chan<- tgbotapi.Chattable {
	return c.sendChan
}

func (c *controller) startSender() {
	for msg := range c.sendChan {
		c.bot.Send(msg)
	}
}
