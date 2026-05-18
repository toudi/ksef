package purchase

import (
	"ksef/internal/invoicesdb/annotations"
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

const (
	xpathIssued     = "//Faktura/Fa/P_1"
	xpathIssuerNIP  = "//Faktura/Podmiot1/DaneIdentyfikacyjne/NIP"
	xpathIssuerName = "//Faktura/Podmiot1/DaneIdentyfikacyjne/Nazwa"
)

func ExtractSeller(
	manager *annotations.Annotations,
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	purchase *types.PurchaseItem,
) error {
	purchase.Date = doc.FindElement(xpathIssued).Text()
	purchase.Seller.NIP = doc.FindElement(xpathIssuerNIP).Text()
	purchase.Seller.Name = doc.FindElement(xpathIssuerName).Text()
	return nil
}
