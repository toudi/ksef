package commands

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	inputprocessors "ksef/internal/sei/input_processors"
	"ksef/internal/uploader"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statefulUpload = &cobra.Command{
	Use:   "stateful-upload [input]",
	Short: "Parsuje plik wejścia oraz wysyła nowe faktury do KSeF",
	RunE:  statefulUploadRun,
	Args:  cobra.ExactArgs(1),
}

const (
	flagConfirm     = "confirm"
	flagAutoCorrect = "auto-correct"
)

var conversionParameters inputprocessors.InputProcessorConfig

func init() {
	flags := statefulUpload.Flags()
	flags.StringVarP(&conversionParameters.CSV.Delimiter, "csv.delimiter", "d", ",", "łańcuch znaków rozdzielający pola (tylko dla CSV)")
	flags.StringVarP(&conversionParameters.CSV.EncodingConversionFile, "csv.encoding", "e", "", "użyj pliku z konwersją znaków (tylko dla CSV)")
	flags.StringVarP(&conversionParameters.XLSX.SheetName, "xlsx.sheet", "s", "", "Nazwa skoroszytu do przetworzenia (tylko dla XLSX)")
	flags.StringVarP(&conversionParameters.Generator, "generator", "g", "fa-3_1.0", "nazwa generatora")
	flags.Bool(flagConfirm, false, "wysyłaj automatycznie po skończonej generacji")
	flags.Bool(flagAutoCorrect, false, "automatycznie wystawiaj faktury korygujące")
	config.CorrectionsFlags(flags)

	RootCommand.AddCommand(statefulUpload)
}

func statefulUploadRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()

	upl, err := uploader.New(
		vip,
		uploader.WithConfirm(vip.GetBool(flagConfirm)),
		uploader.WithAutoCorrect(vip.GetBool(flagAutoCorrect)),
	)
	if err != nil {
		return err
	}

	generator, err := sei.SEI_Init(runtime.GetGateway(vip), conversionParameters, sei.WithInvoiceReadyFunc(upl.InvoiceReady))

	if err != nil {
		return err
	}

	upl.SetGenerator(generator)

	if err = generator.ProcessSourceFile(args[0]); err != nil {
		if err == uploader.ErrAutoCorrectDisabled {
			err = errors.New("tryb automatycznego wystawiania korekt nie został włączony. użyj flagi --auto-correct")
		}
		if err == uploader.ErrRecipientChanged {
			err = errors.New("wykryto zmianę numeru NIP nabywcy faktury. należy wygenerować fakturę korygującą oraz nową fakturę.")
		}
		return err
	}

	return upl.Close()
}
