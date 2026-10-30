package downloaderinterface

import (
	"bytes"
	"context"
	"ksef/core/ksef"
	"ksef/internal/client/v2/types/invoices"
)

type InvoiceDownloader interface {
	Download(
		ctx context.Context,
		invoiceReady func(
			subjectType ksef.SubjectType,
			invoice invoices.InvoiceMetadata,
			content bytes.Buffer,
		) error,
	) error

	Close() error
}
