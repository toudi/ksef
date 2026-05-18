package purchase

import (
	"ksef/internal/invoicesdb/annotations"
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

func ExtractRefNos(
	manager *annotations.Annotations,
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	purchase *types.PurchaseItem,
) error {
	purchase.RefNo = invoice.RefNo
	purchase.KSeFRefNo = invoice.KSeFRefNo
	return nil
}
