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
	issuer      string
}

var MetadataCommand *metadataCommand
var metadataArgs = &metadataArgsType{}

func init() {
	MetadataCommand = &metadataCommand{
		Command: Command{
			Name:        "metadata",
			FlagSet:     flag.NewFlagSet("metadata", flag.ExitOnError),
			Description: "generuje plik metadanych dla wskazanego katalogu faktur",
			Run:         metadataRun,
			Args:        metadataArgs,
		},
	}

	MetadataCommand.FlagSet.BoolVar(&metadataArgs.testGateway, "t", false, "użyj bramki testowej")
	MetadataCommand.FlagSet.StringVar(&metadataArgs.path, "p", "", "ścieżka do wygenerowanych plików")
	MetadataCommand.FlagSet.StringVar(&metadataArgs.issuer, "i", "", "numer NIP wystawcy faktur")

	registerCommand(&MetadataCommand.Command)
}

func metadataRun(c *Command) error {
	if metadataArgs.path == "" || metadataArgs.issuer == "" {
		c.FlagSet.Usage()
		return nil
	}

	fmt.Printf("generowanie metadanych\n")

	meta := &metadata.Metadata{CertificateFile: "klucze/prod/publicKey.pem", Issuer: metadataArgs.issuer}
	if metadataArgs.testGateway {
		meta.CertificateFile = "klucze/test/publicKey.pem"
	}

	return meta.Prepare(metadataArgs.path)
}
