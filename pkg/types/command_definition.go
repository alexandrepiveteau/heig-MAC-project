package types

// CmdInstantiation a type representing a way to start a bot command
type CmdInstantiation func(chan interface{}) Comm

// A CommandDefinition is used to share available commands, how we can call
// them and what they are used for
type CommandDefinition struct {
	// The name of the command
	Command string
	// A description of the command
	Description string
	// A way to start the command
	Instantiation CmdInstantiation
}
