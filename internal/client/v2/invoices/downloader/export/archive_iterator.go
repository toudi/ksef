package export

import (
	"bytes"
	"ksef/internal/client/v2/types/invoices"
)

type archivedInvoice struct {
	Metadata      invoices.InvoiceMetadata
	ContentBuffer bytes.Buffer
}
