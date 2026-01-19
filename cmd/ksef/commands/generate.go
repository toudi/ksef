package commands

import (
	"bytes"
	"errors"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei"
	inputprocessors "ksef/internal/sei/input_processors"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errTargetDirectoryDoesNotExist = errors.New("katalog docelowy nie istnieje. użyj flagi --mkdir / -m jeśli chcesz go stworzyć")

var generateCommand = &cobra.Command{
	Use:   "generate [input]",
	Short: "konwertuje plik CSV/YAML/XLSX do pliku KSeF (XML)",
	RunE:  generateRun,
	Args:  cobra.ExactArgs(1),
}

func init() {
	flags := generateCommand.Flags()
	inputprocessors.GeneratorFlags(flags)
	flags.StringP("output", "o", "", "nazwa katalogu wyjściowego")
	flags.BoolP("mkdir", "m", false, "stwórz katalog, jeśli nie istnieje")

	RootCommand.AddCommand(generateCommand)
}

func generateRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	output := vip.GetString("output")

	logging.GenerateLogger.Info("generate")
	fileName := args[0]
	if fileName == "" || output == "" {
		return cmd.Usage()
	}

	// check if the registry directory does not exist:
	if _, err := os.Stat(output); os.IsNotExist(err) {
		if !viper.GetBool("mkdir") {
			return errTargetDirectoryDoesNotExist
		}
		if err = os.MkdirAll(output, 0775); err != nil {
			return err
		}
	}

	var xmlBuffer bytes.Buffer
	var invoiceOrdNo int

	sei, err := sei.SEI_Init(
		vip,

		sei.WithInvoiceReadyFunc(func(i *sei.ParsedInvoice) error {
			xmlBuffer.Reset()
			outputFile, err := os.Create(path.Join(output, fmt.Sprintf("invoice-%d.xml", invoiceOrdNo)))
			if err != nil {
				return err
			}
			defer outputFile.Close()
			if err = i.ToXML(time.Time{}, outputFile); err != nil {
				return err
			}
			invoiceOrdNo += 1
			return nil
		}),
	)
	if err != nil {
		return err
	}
	if err = sei.ProcessSourceFile(fileName); err != nil {
		return fmt.Errorf("error calling processSourceFile: %v", err)
	}

	return nil
}
