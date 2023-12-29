package commands

import (
	"flag"
	"fmt"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/session/interactive"
)

type saveTokenCommand struct {
	Command
}

type saveTokenArgsType struct {
	token       string
	NIP         string
	testGateway bool
}

var SaveTokenCommand *saveTokenCommand
var saveTokenArgs = &saveTokenArgsType{}

func init() {
	SaveTokenCommand = &saveTokenCommand{
		Command: Command{
			Name:        "save-token",
			FlagSet:     flag.NewFlagSet("save-token", flag.ExitOnError),
			Description: "zapisuje token KSeF w systemowym pęku kluczy",
			Run:         saveTokenRun,
			Args:        saveTokenArgs,
		},
	}

	SaveTokenCommand.FlagSet.BoolVar(&saveTokenArgs.testGateway, "t", false, "użyj bramki testowej")
	SaveTokenCommand.FlagSet.StringVar(&saveTokenArgs.NIP, "nip", "", "numer NIP podatnika")
	SaveTokenCommand.FlagSet.StringVar(&saveTokenArgs.token, "token", "", "token wygenerowany na środowisku KSeF")

	registerCommand(&SaveTokenCommand.Command)
}

func saveTokenRun(c *Command) error {
	if saveTokenArgs.token == "" || saveTokenArgs.NIP == "" {
		c.FlagSet.Usage()
		return nil
	}

	var environment = client.ProductionEnvironment
	if saveTokenArgs.testGateway {
		environment = client.TestEnvironment
	}

	gateway, err := client.APIClient_Init(environment)
	if err != nil {
		return fmt.Errorf("unknown environment: %v", environment)
	}

	session := interactive.InteractiveSessionInit(gateway)
	return session.PersistToken(saveTokenArgs.NIP, saveTokenArgs.token)
}
