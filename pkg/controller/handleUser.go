package controller

import (
	"climb/pkg/types"
	"climb/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleUser(
	ctrl Controller,
	updates <-chan tgbotapi.Update,
) {
	var forwarder *types.Comm
	commandTermination := make(chan interface{})

	for {
		select {
		case <-commandTermination:
			// commands wants to end
			if forwarder != nil {
				forwarder.StopCommand <- struct{}{}
				forwarder = nil
			}
			break

		case update := <-updates:
			// we received an update message from our user
			utils.LogReception(update)

			if update.Message != nil && update.Message.IsCommand() {

				// Clean up previous commands
				if forwarder != nil {
					forwarder.StopCommand <- struct{}{}
					forwarder = nil
				}

				// Get new command started
				for _, cmd := range ctrl.AvailableCommands() {
					if update.Message.Command() == cmd.Command {
						comm := cmd.Instantiation(commandTermination)
						forwarder = &comm
						forwarder.Updates <- update
						break
					}
				}
			} else if forwarder != nil {
				forwarder.Updates <- update
			}
			break
		}
	}
}
