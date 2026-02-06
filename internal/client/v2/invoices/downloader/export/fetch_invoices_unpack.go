package export

import (
	"bytes"
	"context"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/encryption"
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
	archiveHandler, err := NewExportArchiveHandler(cipher)
	if err != nil {
		return err
	}
	if err = archiveHandler.DownloadExportFile(ctx, statusResponse); err != nil {
		return err
	}

	var invoiceContent bytes.Buffer

	for _, invoice := range archiveHandler.contents.Invoices {
		if err = archiveHandler.ReadInvoice(invoice.KSeFNumber, &invoiceContent); err != nil {
			return err
		}

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
