package commands

import (
	"context"
	"errors"
	"flag"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type CommandCallable = func(c *Command) error

type Command struct {
	Name        string
	FlagSet     *flag.FlagSet
	Description string
	Run         CommandCallable
	Context     context.Context
}

type commandsRegistry []*Command

func (r *commandsRegistry) GetByName(name string) *Command {
	for _, command := range *r {
		if command.Name == name {
			return command
		}
	}

	return nil
}

var Registry commandsRegistry
var MaxCommandName int

func registerCommand(command *Command) {
	if Registry == nil {
		Registry = make(commandsRegistry, 0)
	}

	Registry = append(Registry, command)
	if len(command.Name) > MaxCommandName {
		MaxCommandName = len(command.Name)
	}
}
