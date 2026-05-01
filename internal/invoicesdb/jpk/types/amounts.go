package types

import (
	"fmt"
	"ksef/internal/invoicesdb/jpk/abstract/processors/vat"
)

type VATAmounts struct {
	ByRate map[VatRate]*VATInfo
	Total  *VATInfo
}

func VATAmounts_Init() *VATAmounts {
	return &VATAmounts{
		Total:  &VATInfo{},
		ByRate: make(map[VatRate]*VATInfo),
	}
}

func (v *VATAmounts) Add(vatInfo *vat.VatInfo) error {
	// first, lookup VAT Rate by string
	vatRate, exists := VatRates[vatInfo.VatRate]
	if !exists {
		return fmt.Errorf("unrecognizable VAT Rate: %s", vatInfo.VatRate)
	}

	vi := VATInfo{
		Base: vatInfo.NetAmount,
		Vat:  vatInfo.VatAmount,
	}

	v.Total.Add(vi)

	if _, exists := v.ByRate[vatRate]; !exists {
		v.ByRate[vatRate] = &VATInfo{}
	}

	v.ByRate[vatRate].Add(vi)

	return nil
}
