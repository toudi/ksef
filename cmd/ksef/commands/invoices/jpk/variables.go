package jpk

import (
	"ksef/internal/invoicesdb/jpk"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

var (
	jpkManager      *jpk.JPKManager
	invoiceChecksum string
	invoice         *monthlyregistry.Invoice
	xmlInvoice      *monthlyregistry.XMLInvoice
	invoiceRegistry *monthlyregistry.Registry
)
