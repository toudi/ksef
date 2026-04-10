package types

import "ksef/internal/money"

type VATInfo struct {
	Base money.MonetaryValue // podstawa
	Vat  money.MonetaryValue // należny podatek
}

func (vi *VATInfo) Add(other VATInfo) {
	vi.Base = vi.Base.Add(other.Base)
	vi.Vat = vi.Vat.Add(other.Vat)
}
