package generators

import (
	"fmt"
	"ksef/common"
	"ksef/common/xml"
)

type FieldToVATRatesMapping map[string][]string

func populateTotalAmounts(root *xml.Node, invoice *common.Invoice, mapping FieldToVATRatesMapping, aggregator func(common.Amount) float64) {
	var totalPerField map[string]float64 = make(map[string]float64)

	for field, vatRatesList := range mapping {
		for _, rate := range vatRatesList {
			totalPerVatRate, exists := invoice.TotalPerVATRate[rate]
			if exists {
				totalPerFieldEntry, exists := totalPerField[field]
				if !exists {
					totalPerFieldEntry = 0
				}
				totalPerFieldEntry += aggregator(totalPerVatRate)
				totalPerField[field] = totalPerFieldEntry
			}
		}
		if totalPerFieldEntry, exists := totalPerField[field]; exists {
			root.SetValue("Faktura.Fa."+field, fmt.Sprintf("%.2f", totalPerFieldEntry/100))
		}
	}
}

func populateTotalVATAmounts(root *xml.Node, invoice *common.Invoice, mapping FieldToVATRatesMapping) {
	populateTotalAmounts(root, invoice, mapping, func(totalItem common.Amount) float64 { return float64(totalItem.VAT) })
}

func populateTotalNetAmounts(root *xml.Node, invoice *common.Invoice, mapping FieldToVATRatesMapping) {
	populateTotalAmounts(root, invoice, mapping, func(totalItem common.Amount) float64 { return float64(totalItem.Net) })
}
