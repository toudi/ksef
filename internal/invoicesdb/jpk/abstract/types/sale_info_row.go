package types

import (
	"ksef/internal/invoicesdb/jpk/abstract/processors/vat"
	"ksef/internal/invoicesdb/jpk/types"
)

type Buyer struct {
	NIP  string
	Name string
}

type SaleItem struct {
	sales      *Sales
	RefNo      string // internal invoice number
	KSeFRefNo  string // ksef invoice number
	VATAmounts *types.VATAmounts
	IssueDate  string
	SaleDate   string
	Buyer      Buyer
}

func (sr *SaleItem) AddAmount(vatInfo *vat.VatInfo) error {
	if err := sr.VATAmounts.Add(vatInfo); err != nil {
		return err
	}

	if err := sr.sales.VATAmounts.Add(vatInfo); err != nil {
		return err
	}

	return nil
}
