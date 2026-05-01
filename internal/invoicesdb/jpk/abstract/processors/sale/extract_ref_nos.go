package sale

import (
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

func ExtractRefNos(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	salesRow *types.SaleItem,
) error {
	salesRow.RefNo = invoice.RefNo
	salesRow.KSeFRefNo = invoice.KSeFRefNo
	return nil
}
