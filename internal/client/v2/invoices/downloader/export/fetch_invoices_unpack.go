package export

import (
	"bytes"
	"context"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/encryption"
	"ksef/internal/logging"
)

func (ed *exportDownloader) downloadAndExtract(
	ctx context.Context,
	cipher *encryption.Cipher,
	exportRequest exportRequest,
	statusResponse *exportStatusResponse,
	invoiceReady func(
		subjectType invoices.SubjectType,
		invoice invoices.InvoiceMetadata,
		content bytes.Buffer,
	) error,
) error {
	var err error
	ed.archiveHandler, err = NewExportArchiveHandler(cipher)
	if err != nil {
		return err
	}
	logging.DownloadLogger.Debug("archiveHandler::DownloadExportFile")
	if err = ed.archiveHandler.DownloadExportFile(ctx, statusResponse); err != nil {
		return err
	}

	var invoiceContent bytes.Buffer
	logging.DownloadLogger.Debug("Liczba faktur w archiwum", "count", len(ed.archiveHandler.contents.Invoices))

	for _, invoice := range ed.archiveHandler.contents.Invoices {
		logging.DownloadLogger.Debug("Odczytuję fakturę z archiwum", "KSeFRefNo", invoice.KSeFNumber)
		if err = ed.archiveHandler.ReadInvoice(invoice.KSeFNumber, &invoiceContent); err != nil {
			logging.DownloadLogger.Error("Błąd odczytu faktury", "err", err)
			return err
		}

		logging.DownloadLogger.Debug("call invoiceReady")
		if err = invoiceReady(
			exportRequest.Filters.SubjectType,
			invoice,
			invoiceContent,
		); err != nil {
			return err
		}
	}

	return nil
}
