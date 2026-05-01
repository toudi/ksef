package sale

import (
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

const (
	xpathIssued = "//Faktura/Fa/P_1"
	xpathSale   = "//Faktura/Fa/P_6"
)

func ExtractDates(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	salesRow *types.SaleItem,
) error {
	salesRow.IssueDate = doc.FindElement(xpathIssued).Text()
	salesRow.SaleDate = doc.FindElement(xpathSale).Text()
	return nil
}
