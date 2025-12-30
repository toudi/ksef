package invoicesdb

import (
	"context"
	"errors"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	"ksef/internal/pdf/printer"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	errUnableToDetectNIP           = errors.New("unable to detect NIP based on imported invoice files")
	errUnableToSaveInvoiceRegistry = errors.New("unable to save monthly registry")
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
	var nip string

	if len(idb.newInvoices) == 0 {
		logging.InvoicesDBLogger.Info("brak nowo zaimportowanych faktur")
		return nil
	}

	if confirm {
		affectedRegistries := make(map[*monthlyregistry.Registry]bool)

		var printer printer.PDFPrinter

		printer, err = pdf.GetInvoicePrinter(vip, "invoice:issued")
		if err != nil {
			logging.PDFRendererLogger.Error("błąd inicjalizacji silnika PDF", "err", err)
		}

		for _, newInvoice := range idb.newInvoices {
			affectedRegistries[newInvoice.registry] = true
			if nip == "" {
				nip = newInvoice.registry.GetNIP()
			}
			if !newInvoice.invoice.Offline {
				// if the invoice isn't issued in the offline mode we cannot print it (yet) - we
				// have to wait for the reference no obtained from KSeF
				continue
			}
			// let's try to print the offline invoice - though only if the engine has been selected
			if printer != nil {
				invoiceFilename := newInvoice.registry.InvoiceFilename(newInvoice.invoice)
				if err = printer.PrintInvoice(
					invoiceFilename.XML,
					invoiceFilename.PDF,
					newInvoice.invoice.GetPrintingMeta(),
				); err != nil {
					return err
				}
			} else {
				logging.InvoicesDBLogger.Warn("wybrano opcję wydruku faktury do PDF ale nie udało się zainicjować silnika wydruku - pomijam wydruk faktury")
			}
		}

		if err = idb.Save(); err != nil {
			logging.InvoicesDBLogger.Error("błąd zapisu rejestru faktur", "err", err)
			return err
		}

		// save the registries so that we persist changes in ord nums
		for registry := range affectedRegistries {
			if err = registry.Save(); err != nil {
				return errors.Join(errUnableToSaveInvoiceRegistry, err)
			}
		}
	}

	if !importConfig.Upload.Enabled {
		return nil
	}

	if nip == "" {
		return errors.Join(errUnableToDetectNIP)
	}

	runtime.SetNIP(idb.vip, nip)

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
