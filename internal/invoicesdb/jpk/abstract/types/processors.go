package types

import (
	"ksef/internal/invoicesdb/annotations"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

type SaleInvoiceProcessor func(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	sale *SaleItem,
) error

type PurchaseInvoiceProcessor func(
	manager *annotations.Annotations,
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	purchase *PurchaseItem,
) error
