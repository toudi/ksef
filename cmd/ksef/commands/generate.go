package commands

import (
	"errors"
	"fmt"
	"ksef/internal/config"
	"ksef/internal/logging"
	"ksef/internal/sei"
	inputprocessors "ksef/internal/sei/input_processors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errTargetDirectoryDoesNotExist = errors.New("katalog docelowy nie istnieje. użyj flagi --mkdir / -m jeśli chcesz go stworzyć")

var generateCommand = &cobra.Command{
	Use:   "generate [input]",
	Short: "konwertuje plik CSV/YAML/XLSX do pliku KSeF (XML) i tworzy rejestr gotowy do wysłania",
	RunE:  generateRun,
	Args:  cobra.ExactArgs(1),
}

type generateArgsType struct {
	Output                 string
	Delimiter              string
	GeneratorName          string
	EncodingConversionFile string
	SheetName              string
	Offline                bool
}

var generateArgs = &generateArgsType{}

func init() {
	flags := generateCommand.Flags()
	flags.StringVarP(&generateArgs.Output, "output", "o", "", "nazwa katalogu wyjściowego")
	flags.StringVarP(&generateArgs.Delimiter, "csv.delimiter", "d", ",", "łańcuch znaków rozdzielający pola (tylko dla CSV)")
	flags.StringVarP(&generateArgs.EncodingConversionFile, "csv.encoding", "e", "", "użyj pliku z konwersją znaków (tylko dla CSV)")
	flags.StringVarP(&generateArgs.SheetName, "xlsx.sheet", "s", "", "Nazwa skoroszytu do przetworzenia (tylko dla XLSX)")
	flags.StringVarP(&generateArgs.GeneratorName, "generator", "g", "fa-3-1", "nazwa generatora")
	flags.BoolVar(&generateArgs.Offline, "offline", false, "oznacz faktury jako generowane w trybie offline")
	flags.BoolP("mkdir", "m", false, "stwórz katalog rejestru, jeśli nie istnieje")

	RootCommand.AddCommand(generateCommand)
}

func generateRun(cmd *cobra.Command, args []string) error {
	logging.GenerateLogger.Info("generate")
	fileName := args[0]
	if fileName == "" || generateArgs.Output == "" {
		return cmd.Usage()
	}

	// check if the registry directory does not exist:
	if _, err := os.Stat(generateArgs.Output); os.IsNotExist(err) {
		if !viper.GetBool("mkdir") {
			return errTargetDirectoryDoesNotExist
		}
		if err = os.MkdirAll(generateArgs.Output, 0775); err != nil {
			return err
		}
	}

	var conversionParameters inputprocessors.InputProcessorConfig
	conversionParameters.CSV.Delimiter = generateArgs.Delimiter
	conversionParameters.CSV.EncodingConversionFile = generateArgs.EncodingConversionFile
	conversionParameters.XLSX.SheetName = generateArgs.SheetName
	conversionParameters.Generator = generateArgs.GeneratorName
	conversionParameters.OfflineMode = generateArgs.Offline

	sei, err := sei.SEI_Init(config.GetGateway(viper.GetViper()), generateArgs.Output, conversionParameters)
	if err != nil {
		return err
	}
	if err = sei.ProcessSourceFile(fileName); err != nil {
		return fmt.Errorf("error calling processSourceFile: %v", err)
	}

	return nil
}
