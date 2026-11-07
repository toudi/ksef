package invoices

import (
	"ksef/cmd/ksef/flags"
	"ksef/core/invoicedb/invoice"
	"ksef/internal/logging"
	"ksef/internal/resources"
	"ksef/internal/runtime"
	"os"

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
	flags.NIP(importPdfCommand.Flags())
}

func importPDFRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	ctx := cmd.Context()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}
	NIP, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	invoiceDataStorage := resources.InvoiceDataStorage(vip)

	invoiceDB, err := resources.InvoiceDB(vip, NIP)
	if err != nil {
		return err
	}

	defer invoiceDB.Close(ctx)

	refNo := args[0]
	storedInvoice, err := invoiceDB.GetInvoiceByRefNo(ctx, refNo)
	if err != nil {
		return err
	}

	storageId := storedInvoice.StorageId(invoice.InvoiceFilePDF)
	srcPDF, err := os.Open(args[1])
	if err != nil {
		return err
	}
	defer srcPDF.Close()

	stored, err := invoiceDataStorage.Store(ctx, srcPDF, storageId, vip.GetBool(flagNameForce))
	if err != nil {
		return err
	}

	if !stored {
		logging.SeiLogger.Info("plik PDF istnieje - pominięto.")
	}

	return nil
}
