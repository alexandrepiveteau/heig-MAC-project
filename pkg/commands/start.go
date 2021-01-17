package commands

import (
	"climb/pkg/types"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Entrypoint of bot command
func StartCmd(
	comm types.Comm,
	commandTermination chan interface{},
	bot *tgbotapi.BotAPI,
	availableCommands []types.CommandDefinition,
) {
	for {
		select {
		case <-comm.StopCommand:
			return

		case update := <-comm.Updates:
			var text string

			for _, cmd := range availableCommands {
				text += fmt.Sprintf(
					"/%s : %s\n",
					cmd.Command,
					cmd.Description,
				)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(msg)

			commandTermination <- struct{}{} // Inform that we have terminated
		}
	}
}
