package invoicesdb

import (
	"context"
	"ksef/internal/invoicesdb/config"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	"ksef/internal/sei"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type invoiceReadyHandlerFunc func(i *sei.ParsedInvoice) error

func (idb *InvoicesDB) Import(
	ctx context.Context,
	vip *viper.Viper,
	srcFilename string,
	confirm bool,
) (err error) {
	importConfig := config.GetImportConfig(vip)

	// initialize importer
	var invoiceReadyHandler invoiceReadyHandlerFunc = dummyInvoiceReadyHandler

	if confirm {
		invoiceReadyHandler = idb.invoiceReady
	} else {
		logging.InvoicesDBLogger.Info("nie wybrano flagi --confirm - żadne dane nie zostaną zapisane na dysku")
	}

	// check if the source filename is of type XML. If yes - then the user is trying to import
	// invoices that are already in FA format.
	if isXML(srcFilename) {
		if err = idb.importXMLInvoices(srcFilename); err != nil {
			return err
		}
	} else {
		// importing from other invoice sources (i.e. csv / yaml / ...)
		importer, err := sei.SEI_Init(vip, sei.WithInvoiceReadyFunc(invoiceReadyHandler))
		if err != nil {
			return err
		}
		if importErr := importer.ProcessSourceFile(srcFilename); importErr != nil {
			return importErr
		}
	}

	// save the database before uploading
	if confirm {
		printer, err := pdf.GetInvoicePrinter(vip, "invoice:issued")
		if err != nil {
			logging.PDFRendererLogger.Error("błąd inicjalizacji silnika PDF", "err", err)
		} else {
			for _, offlineInvoice := range idb.offlineInvoices {
				invoiceFilename := offlineInvoice.registry.InvoiceFilename(offlineInvoice.invoice)
				if err = printer.PrintInvoice(
					invoiceFilename,
					strings.Replace(invoiceFilename, ".xml", ".pdf", 1),
					offlineInvoice.invoice.GetPrintingMeta(),
				); err != nil {
					return err
				}
			}
		}
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
		uploaderconfig.GetUploaderConfig(vip),
		statuscheckerconfig.GetStatusCheckerConfig(vip),
	)
}

func dummyInvoiceReadyHandler(i *sei.ParsedInvoice) error {
	// dummy implementation that does nothing.

	return nil
}

func isXML(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".xml"
}
