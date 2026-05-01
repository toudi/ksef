package purchase

import (
	"ksef/internal/invoicesdb/jpk/abstract/types"
	"ksef/internal/invoicesdb/jpk/manager"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

func ExtractRefNos(
	manager *manager.JPKManager,
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	purchase *types.PurchaseItem,
) error {
	purchase.RefNo = invoice.RefNo
	purchase.KSeFRefNo = invoice.KSeFRefNo
	return nil
}
