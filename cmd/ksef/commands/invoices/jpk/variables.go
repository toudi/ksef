package jpk

import (
	"ksef/internal/invoicesdb/jpk/manager"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

var (
	jpkManager      *manager.JPKManager
	invoiceChecksum string
	invoice         *monthlyregistry.Invoice
	xmlInvoice      *monthlyregistry.XMLInvoice
	invoiceRegistry *monthlyregistry.Registry
)
