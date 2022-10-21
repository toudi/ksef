package commands

import (
	"flag"
	"fmt"
	"ksef/generators"
)

type generateCommand struct {
	Command
}

type generateArgsType struct {
	FileName      string
	Output        string
	Delimiter     string
	GeneratorName string
}

var GenerateCmd *generateCommand
var generateArgs = &generateArgsType{}

func init() {
	GenerateCmd = &generateCommand{
		Command: Command{
			Name:        "generate",
			FlagSet:     flag.NewFlagSet("generate", flag.ExitOnError),
			Description: "Konwertuje plik CSV do pliku KSEF (XML)",
			Run:         generateRun,
			Args:        generateArgs,
		},
	}

	GenerateCmd.FlagSet.StringVar(&generateArgs.FileName, "f", "", "nazwa pliku do przetworzenia")
	GenerateCmd.FlagSet.StringVar(&generateArgs.Output, "o", "", "nazwa pliku wyjściowego")
	GenerateCmd.FlagSet.StringVar(&generateArgs.Delimiter, "d", ",", "łańcuch znaków rozdzielający pola w CSV")
	GenerateCmd.FlagSet.StringVar(&generateArgs.GeneratorName, "g", "fa_1_1", "nazwa generatora")
	GenerateCmd.FlagSet.BoolVar(&metadataArgs.testGateway, "t", false, "użyj certyfikatu bramki testowej do generowania podpisu")

	registerCommand(&GenerateCmd.Command)
}

var _generator *generators.Generator

func generateRun(c *Command) error {
	if generateArgs.FileName == "" || generateArgs.Output == "" {
		GenerateCmd.FlagSet.Usage()
		return nil
	}

	generator, err := generators.Run(generateArgs.GeneratorName, generateArgs.Delimiter, generateArgs.FileName, generateArgs.Output)
	if err != nil {
		return fmt.Errorf("błąd generowania danych wejściowych: %v", err)
	}

	metadataArgs.path = generateArgs.Output
	metadataArgs.issuer = generator.Issuer()
	return metadataRun(c)
}
