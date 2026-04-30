package types

import (
	"ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

type Sales struct {
	Rows       []*SaleItem
	VATAmounts *types.VATAmounts
	processors []SaleInvoiceProcessor
}

func NewSales(processors []SaleInvoiceProcessor) *Sales {
	return &Sales{
		processors: processors,
		VATAmounts: types.VATAmounts_Init(),
	}
}

func (s *Sales) ProcessInvoice(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
) error {
	row := &SaleItem{
		sales:      s,
		VATAmounts: types.VATAmounts_Init(),
	}

	for _, processor := range s.processors {
		if err := processor(invoice, doc, row); err != nil {
			return err
		}
	}

	s.Rows = append(s.Rows, row)

	return nil
}
