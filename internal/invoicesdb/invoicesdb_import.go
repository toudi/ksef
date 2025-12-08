package invoicesdb

import (
	"context"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/logging"
	"ksef/internal/sei"

	"github.com/spf13/viper"
)

type invoiceReadyHandlerFunc func(i *sei.ParsedInvoice) error

func (idb *InvoicesDB) Import(
	ctx context.Context,
	vip *viper.Viper,
	srcFilename string,
	confirm bool,
) error {
	importConfig := config.GetImportConfig(vip)

	// initialize importer
	var invoiceReadyHandler invoiceReadyHandlerFunc = dummyInvoiceReadyHandler

	if confirm {
		invoiceReadyHandler = idb.invoiceReady
	} else {
		logging.InvoicesDBLogger.Info("nie wybrano flagi --confirm - żadne dane nie zostaną zapisane na dysku")
	}

	importer, err := sei.SEI_Init(vip, sei.WithInvoiceReadyFunc(invoiceReadyHandler))
	if err != nil {
		return err
	}
	if importErr := importer.ProcessSourceFile(srcFilename); importErr != nil {
		return importErr
	}

	// save the database before uploading
	if confirm {
		if err := idb.Save(); err != nil {
			logging.InvoicesDBLogger.Error("błąd zapisu rejestru faktur", "err", err)
			return err
		}
	}

	if !importConfig.Upload.Enabled {
		return nil
	}

	// user wants to automatically upload invoices
	return idb.UploadOutstandingInvoices(
		ctx,
		vip,
	)
}

func dummyInvoiceReadyHandler(i *sei.ParsedInvoice) error {
	// dummy implementation that does nothing.

	return nil
}
