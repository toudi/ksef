package downloaderinterface

import (
	"bytes"
	"context"
	"ksef/internal/client/v2/types/invoices"
)

type InvoiceDownloader interface {
	Download(
		ctx context.Context,
		invoiceReady func(
			subjectType invoices.SubjectType,
			invoice invoices.InvoiceMetadata,
			content bytes.Buffer,
		) error,
	) (err error)

	Close() error
}
