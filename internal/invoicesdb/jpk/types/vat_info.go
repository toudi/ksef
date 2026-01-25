package types

import "ksef/internal/money"

type VATInfo struct {
	Base money.MonetaryValue // podstawa
	Vat  money.MonetaryValue // należny podatek
}
