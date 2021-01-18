package controller

import (
	"climb/pkg/types"
	"climb/pkg/utils"
	"log"
)

func handleUser(
	ctrl Controller,
	userdata types.UserData,
) {
	// Add user in neo4j
	err := userdata.RegisterInNeo4j(ctrl.Neo4j())
	if err != nil {
		log.Println(err.Error())
	}

	// Prepare needed controllers
	var forwarder *types.Comm
	commandTermination := make(chan interface{})

	// Main goroutine loop
	for {
		select {
		case <-commandTermination:
			// commands wants to end
			if forwarder != nil {
				forwarder.StopCommand <- struct{}{}
				forwarder = nil
			}
			break

		case update := <-userdata.Channel:
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
						comm := cmd.Instantiation(commandTermination, userdata, ctrl.GetCurrentUsers())
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
