package fa_2

import (
	"ksef/internal/invoice"
)

var fieldToVatRatesMapping *invoice.FieldToVATRatesMapping = invoice.NewFieldToVATRatesMapping(
	map[string][]string{
		"P_13_1": {"22", "23"},
		"P_13_2": {"8", "7"},
		"P_13_3": {"5"},
	},
	map[string][]string{
		"P_14_1": {"22", "23"},
		"P_14_2": {"8", "7"},
		"P_14_3": {"5"},
	},
)
