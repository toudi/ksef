package export

import (
	"bytes"

	"ksef/core/ksef"
)

type archivedInvoice struct {
	Metadata      ksef.InvoiceMetadata
	ContentBuffer bytes.Buffer
}
