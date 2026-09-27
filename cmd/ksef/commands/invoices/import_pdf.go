package invoices

import (
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"

	"github.com/goforj/godump"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importPdfCommand = &cobra.Command{
	Use:   "import-pdf [ref-no] [pdf-file]",
	Short: "importuj zewnętrzny PDF z wydrukiem faktury",
	Args:  cobra.ExactArgs(2),
	RunE:  importPDFRun,
}

const (
	flagNameForce = "force"
)

func init() {
	InvoicesCommand.AddCommand(importPdfCommand)
	importPdfCommand.Flags().BoolP(flagNameForce, "f", false, "nadpisz plik PDF nawet jeśli istnieje")
}

func importPDFRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	refNo := args[0]
	invoice, reg, err := findInvoiceByRefNo(vip, refNo)
	if err != nil {
		return err
	}

	godump.Dump(invoice, reg)

	targetPDFName := reg.InvoiceFilename(invoice).PDF

	if _, err = os.Stat(targetPDFName); err == nil {
		// file exists
		if !vip.GetBool(flagNameForce) {
			logging.SeiLogger.Info("plik PDF istnieje - pomijam.")
			return nil
		}
	}

	return utils.CopyFile(args[1], targetPDFName)
}
