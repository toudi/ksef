package commands

import (
	"flag"
	"fmt"
	"ksef/metadata"
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

	meta := &metadata.Metadata{CertificateFile: "klucze/prod/publicKey.pem"}
	if metadataArgs.testGateway {
		meta.CertificateFile = "klucze/test/publicKey.pem"
	}

	return meta.Prepare(metadataArgs.path)
}
