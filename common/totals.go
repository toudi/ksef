package common

import (
	"fmt"
	"ksef/common/xml"
)

const (
	AccumulatorNET = iota
	AccumulatorVAT = iota
)

// this is the internal type that will represent a reverse mapping of
// MappingNet/MappingVAT structures.
type Aggregator struct {
	Net string
	VAT string
}

type FieldToVATRatesMapping struct {
	MappingNet     map[string][]string
	MappingVAT     map[string][]string
	Totals         map[string]float64
	reverseMapping map[string]*Aggregator
}

func NewFieldToVATRatesMapping(MappingNet map[string][]string, MappingVAT map[string][]string) *FieldToVATRatesMapping {
	mapping := &FieldToVATRatesMapping{
		MappingNet:     MappingNet,
		MappingVAT:     MappingVAT,
		Totals:         make(map[string]float64),
		reverseMapping: make(map[string]*Aggregator),
	}

	for outputField, vatRatesList := range MappingNet {
		for _, rate := range vatRatesList {
			mapping.reverseMapping[rate] = &Aggregator{Net: outputField}
		}
	}

	for outputField, vatRatesList := range MappingVAT {
		for _, rate := range vatRatesList {
			mapping.reverseMapping[rate].VAT = outputField
		}
	}

	return mapping
}

func (f *FieldToVATRatesMapping) Accumulate(item *InvoiceItem) {
	vatRate := item.UnitPrice.Vat.Description

	reverseMapping, exists := f.reverseMapping[vatRate]
	if !exists {
		return
	}

	if _, exists := f.Totals[reverseMapping.Net]; !exists {
		f.Totals[reverseMapping.Net] = 0.0
	}
	if _, exists := f.Totals[reverseMapping.VAT]; !exists {
		f.Totals[reverseMapping.VAT] = 0.0
	}

	f.Totals[reverseMapping.Net] += float64(item.Amount().Net)
	f.Totals[reverseMapping.VAT] += float64(item.Amount().VAT)
}

func (f *FieldToVATRatesMapping) Populate(root *xml.Node) {
	for totalFieldName, totalValue := range f.Totals {
		root.SetValue("Faktura.Fa."+totalFieldName, fmt.Sprintf("%.2f", totalValue/100))
	}
}

func (f *FieldToVATRatesMapping) Zero() {
	f.Totals = make(map[string]float64)
}
