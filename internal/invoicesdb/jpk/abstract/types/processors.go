package types

import (
	"ksef/internal/invoicesdb/jpk/manager"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

type SaleInvoiceProcessor func(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	sale *SaleItem,
) error

type PurchaseInvoiceProcessor func(
	manager *manager.JPKManager,
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	purchase *PurchaseItem,
) error
