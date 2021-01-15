package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func Color(
	updates <-chan tgbotapi.Update,
	send chan<- tgbotapi.Chattable,
) {
}
