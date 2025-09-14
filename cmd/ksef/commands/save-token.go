package commands

import (
	"flag"
	"ksef/internal/config"
	kseftoken "ksef/internal/sei/api/client/v2/auth/ksef_token"
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

	var env config.APIEnvironment = config.APIEnvironmentProd
	if saveTokenArgs.testGateway {
		env = config.APIEnvironmentTest
	}

	return kseftoken.PersistKsefTokenToKeyring(config.GetConfig().APIConfig(env).Host, saveTokenArgs.NIP, saveTokenArgs.token)
}
