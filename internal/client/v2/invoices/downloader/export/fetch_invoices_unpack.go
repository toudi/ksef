package export

import (
	"bytes"
	"context"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/encryption"
)

func (ed *ExportDownloader) downloadAndExtract(
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
	ed.archiveHandler, err = NewExportArchiveHandler(cipher, ed.logger)
	if err != nil {
		return err
	}
	ed.logger.Debug("archiveHandler::DownloadExportFile")
	if err = ed.archiveHandler.DownloadExportFile(ctx, statusResponse); err != nil {
		return err
	}

	var invoiceContent bytes.Buffer
	ed.logger.Debug("Liczba faktur w archiwum", "count", len(ed.archiveHandler.contents.Invoices))

	for _, invoice := range ed.archiveHandler.contents.Invoices {
		ed.logger.Debug("Odczytuję fakturę z archiwum", "KSeFRefNo", invoice.KSeFNumber)
		if err = ed.archiveHandler.ReadInvoice(invoice.KSeFNumber, &invoiceContent); err != nil {
			ed.logger.Error("Błąd odczytu faktury", "err", err)
			return err
		}

		ed.logger.Debug("call invoiceReady")
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
