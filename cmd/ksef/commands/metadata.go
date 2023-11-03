package commands

import (
	"flag"
	"fmt"
	"ksef/api"
)

type metadataCommand struct {
	Command
}

type metadataArgsType struct {
	path        string
	testGateway bool
}

var MetadataCommand *metadataCommand
var metadataArgs = &metadataArgsType{}

func init() {
	MetadataCommand = &metadataCommand{
		Command: Command{
			Name:        "metadata",
			FlagSet:     flag.NewFlagSet("metadata", flag.ExitOnError),
			Description: "generuje plik metadanych dla wskazanego katalogu faktur (tylko tryb wsadowy)",
			Run:         metadataRun,
			Args:        metadataArgs,
		},
	}

	MetadataCommand.FlagSet.BoolVar(&metadataArgs.testGateway, "t", false, "użyj bramki testowej")
	MetadataCommand.FlagSet.StringVar(&metadataArgs.path, "p", "", "ścieżka do wygenerowanych plików")

	registerCommand(&MetadataCommand.Command)
}

func metadataRun(c *Command) error {
	if metadataArgs.path == "" {
		c.FlagSet.Usage()
		return nil
	}

	fmt.Printf("generowanie metadanych\n")

	var environment = api.ProductionEnvironment
	if metadataArgs.testGateway {
		environment = api.TestEnvironment
	}

	gateway, err := api.API_Init(environment)
	if err != nil {
		return fmt.Errorf("unknown environment: %d", environment)
	}

	batchSession := gateway.BatchSessionInit()
	return batchSession.GenerateMetadata(metadataArgs.path)
}
