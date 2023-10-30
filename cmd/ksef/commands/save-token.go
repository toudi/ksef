package commands

import (
	"flag"
	"fmt"
	"ksef/common"

	"github.com/zalando/go-keyring"
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

	_ = registerCommand(&SaveTokenCommand.Command)
}

func saveTokenRun(c *Command) error {
	if saveTokenArgs.token == "" || saveTokenArgs.NIP == "" {
		c.FlagSet.Usage()
		return nil
	}

	var gateway = common.KSeFHost
	if saveTokenArgs.testGateway {
		gateway = common.KSeFTestHost
	}

	err := keyring.Set(gateway, saveTokenArgs.NIP, saveTokenArgs.token)
	fmt.Printf("Err for keyring.set: %v\n", err)

	fmt.Printf("try to read the password:\n")
	token, err := keyring.Get(gateway, saveTokenArgs.NIP)
	fmt.Printf("token=%s; err=%v\n", token, err)
	return nil
}
