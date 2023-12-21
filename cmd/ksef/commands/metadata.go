package commands

import (
	"flag"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/session/batch"
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

	logging.SeiLogger.Info().Msg("generowanie metadanych")

	var environment = client.ProductionEnvironment
	if metadataArgs.testGateway {
		environment = client.TestEnvironment
	}

	gateway, err := client.APIClient_Init(environment)
	if err != nil {
		return fmt.Errorf("unknown environment: %d", environment)
	}

	batchSession := batch.BatchSessionInit(gateway)
	return batchSession.GenerateMetadata(metadataArgs.path)
}
