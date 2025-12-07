package invoices

import (
	"ksef/internal/invoicesdb"
	"ksef/internal/logging"
	"ksef/internal/sei"
	inputprocessors "ksef/internal/sei/input_processors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCommand = &cobra.Command{
	Use:   "import [input]",
	Short: "importuj faktury z pliku do bazy",
	Args:  cobra.ExactArgs(1),
	RunE:  importRun,
}

func init() {
	inputprocessors.GeneratorFlags(importCommand.Flags())
	importCommand.Flags().SortFlags = false
	InvoicesCommand.AddCommand(importCommand)
}

func importRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	// initialize the invoicesdb
	invoicesDB, err := invoicesdb.NewInvoicesDB(vip)
	if err != nil {
		return err
	}
	// initialize importer
	var invoiceReadyHandler = func(i *sei.ParsedInvoice) error {
		// dummy implementation that does nothing.
		return nil
	}
	if vip.GetBool(flagNameConfirm) {
		invoiceReadyHandler = invoicesDB.InvoiceReady
		defer func() {
			if err := invoicesDB.Save(); err != nil {
				logging.InvoicesDBLogger.Error("błąd zapisu rejestru faktur", "err", err)
			}
		}()
	} else {
		logging.InvoicesDBLogger.Info("nie wybrano flagi --confirm - żadne dane nie zostaną zapisane na dysku")
	}
	importer, err := sei.SEI_Init(vip, sei.WithInvoiceReadyFunc(invoiceReadyHandler))
	if err != nil {
		return err
	}
	return importer.ProcessSourceFile(args[0])
}
