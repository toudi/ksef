package constants

const (
	RefundMode15Days    = "15-bank"
	RefundMode25DaysVat = "25-vat"
	RefundMode25Days    = "25-bank"
	RefundMode40Days    = "40-bank"
	RefundMode180Days   = "180-bank"
)

var VATRefundModes = []string{
	RefundMode15Days,
	RefundMode25DaysVat,
	RefundMode25Days,
	RefundMode40Days,
	RefundMode180Days,
}

const (
	FlagNameSurplusAction = "jpk.surplus.mode"
	FlagNameRefundMode    = "jpk.surplus.refund"
	FlagNameOffsetTaxCode = "jpk.surplus.offset-tax"
)
