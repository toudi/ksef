package commands

import (
	"flag"
	"fmt"
	"ksef/common"
	"ksef/generators"
)

type generateCommand struct {
	Command
}

type generateArgsType struct {
	FileName               string
	Output                 string
	Delimiter              string
	GeneratorName          string
	EncodingConversionFile string
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
	GenerateCmd.FlagSet.StringVar(&generateArgs.GeneratorName, "g", "fa-1-1", "nazwa generatora")
	GenerateCmd.FlagSet.BoolVar(&metadataArgs.testGateway, "t", false, "użyj certyfikatu bramki testowej do generowania podpisu")
	GenerateCmd.FlagSet.StringVar(&generateArgs.EncodingConversionFile, "e", "", "użyj pliku z konwersją znaków")

	registerCommand(&GenerateCmd.Command)
}

var _generator *common.Generator

func generateRun(c *Command) error {
	if generateArgs.FileName == "" || generateArgs.Output == "" {
		GenerateCmd.FlagSet.Usage()
		return nil
	}

	generator, err := generators.Run(generateArgs.GeneratorName, generateArgs.Delimiter, generateArgs.FileName, generateArgs.Output, generateArgs.EncodingConversionFile)
	if err != nil {
		return fmt.Errorf("błąd generowania danych wejściowych: %v", err)
	}

	metadataArgs.path = generateArgs.Output
	metadataArgs.generator = generateArgs.GeneratorName
	metadataArgs.issuer = generator.IssuerTIN()
	return metadataRun(c)
}
