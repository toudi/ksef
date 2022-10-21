package main

import (
	"fmt"
	"ksef/cmd/ksef/commands"
	"os"
)

var command *commands.Command

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Printf("Please specify at least one sub-command\nAvailable subcommands:\n\n")
		for _, command := range commands.Registry {
			fmt.Printf("%-*s - %s\n", commands.MaxCommandName, command.Name, command.Description)
		}
		return
	}

	command = commands.Registry.GetByName(os.Args[1])
	if command == nil {
		fmt.Printf("unknown command\n")
		return
	}

	command.FlagSet.Parse(os.Args[2:])
	err = command.Run(command)
	if err != nil {
		fmt.Printf("error running %s:\n%s\n", os.Args[1], err)
		return
	}
}
