package commands

import (
	"climb/pkg/comm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// A CommandDescription is used to share available commands, how we can call
// them and what they are used for
type CommandDescription struct {
	Command     string
	Description string
}

var Start = CommandDescription{
	Command:     "start",
	Description: "The start command shows available commands",
}

// Entrypoint of bot command
func StartCmd(
	comm comm.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
) {
	for {
		select {
		case <-comm.StopCommand:
			return

		case update := <-comm.Updates:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I'm lost.")
			bot.Send(msg)

			commandTermination <- struct{}{} // Inform that we have terminated
		}
	}
}
