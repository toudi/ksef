package commands

import (
	"flag"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei"
	inputprocessors "ksef/internal/sei/input_processors"
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
	SheetName              string
}

var GenerateCmd *generateCommand
var generateArgs = &generateArgsType{}

func init() {
	GenerateCmd = &generateCommand{
		Command: Command{
			Name:        "generate",
			FlagSet:     flag.NewFlagSet("generate", flag.ExitOnError),
			Description: "Konwertuje plik CSV/YAML/XLSX do pliku KSEF (XML)",
			Run:         generateRun,
			Args:        generateArgs,
		},
	}

	GenerateCmd.FlagSet.StringVar(&generateArgs.FileName, "f", "", "nazwa pliku do przetworzenia")
	GenerateCmd.FlagSet.StringVar(&generateArgs.Output, "o", "", "nazwa katalogu wyjściowego")
	GenerateCmd.FlagSet.StringVar(
		&generateArgs.Delimiter,
		"d",
		",",
		"łańcuch znaków rozdzielający pola (tylko dla CSV)",
	)
	GenerateCmd.FlagSet.StringVar(
		&generateArgs.SheetName,
		"s",
		"",
		"Nazwa skoroszytu do przetworzenia (tylko dla XLSX)",
	)
	GenerateCmd.FlagSet.StringVar(&generateArgs.GeneratorName, "g", "fa-2", "nazwa generatora")
	GenerateCmd.FlagSet.StringVar(
		&generateArgs.EncodingConversionFile,
		"e",
		"",
		"użyj pliku z konwersją znaków (tylko dla CSV)",
	)

	registerCommand(&GenerateCmd.Command)
}

func generateRun(c *Command) error {
	logging.GenerateLogger.Info("generate")
	if generateArgs.FileName == "" || generateArgs.Output == "" {
		GenerateCmd.FlagSet.Usage()
		return nil
	}

	var conversionParameters inputprocessors.InputProcessorConfig
	conversionParameters.CSV.Delimiter = generateArgs.Delimiter
	conversionParameters.CSV.EncodingConversionFile = generateArgs.EncodingConversionFile
	conversionParameters.XLSX.SheetName = generateArgs.SheetName
	conversionParameters.Generator = generateArgs.GeneratorName
	sei, err := sei.SEI_Init(generateArgs.Output, conversionParameters)
	if err != nil {
		return err
	}
	if err = sei.ProcessSourceFile(generateArgs.FileName); err != nil {
		return fmt.Errorf("error calling processSourceFile: %v", err)
	}

	return nil
}
