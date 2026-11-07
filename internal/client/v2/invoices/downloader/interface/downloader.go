package downloaderinterface

import (
	"bytes"
	"context"
	"ksef/core/downloader"
	"ksef/core/ksef"
)

type InvoiceDownloader interface {
	Download(
		ctx context.Context,
		params downloader.Config,
		invoiceReady func(
			subjectType ksef.SubjectType,
			invoice ksef.InvoiceMetadata,
			content bytes.Buffer,
		) error,
	) error

	Close() error
}
