package fa

import (
	"ksef/internal/invoice"
)

var fieldToVatRatesMapping *invoice.FieldToVATRatesMapping = invoice.NewFieldToVATRatesMapping(
	map[string][]string{
		"P_13_1": {"22", "23"},
		"P_13_2": {"8", "7"},
		"P_13_3": {"5"},
		// KSeF web UI: "np na podstawie art. 100 ust. 1 pkt 4 ustawy"
		"P_13_9": {"np", "np II"},
		// KSeF web UI: "np z wyłączeniem art. 100 ust 1 pkt 4 ustawy"
		"P_13_8": {"np.except", "np I"},
	},
	map[string][]string{
		"P_14_1": {"22", "23"},
		"P_14_2": {"8", "7"},
		"P_14_3": {"5"},
		// there's no need to calculate VAT for these rates as the numerical
		// value will always be 0, however because they exist in the above
		// mapping we have to somehow indicate to totals.go to ignore them.
		invoice.IgnoreRate: {"np", "np I", "np II", "np.except"},
	},
)
