package fa

import (
	"fmt"
	"ksef/internal/invoice"
	"ksef/internal/money"
	"ksef/internal/xml"
	"time"
)

// this function is meant to be used when you have a invoice.Invoice object as your source
// and want to output the xml.Node object as your output, which is then used to generate
// metadata and so on.
// if this is not what you wish to do, then please use the line-based generator approach (either CSV or XLSX)
func (fg *FAGenerator) InvoiceToXMLTree(invoice *invoice.Invoice) (*xml.Node, error) {
	root := &xml.Node{Name: "Faktura"}

	root.SetValuesFromMap(fg.commonData)

	if invoice.Correction != nil {
		delete(invoice.Attributes, "Faktura.Fa.P_6")
	}

	root.SetValuesFromMap(invoice.Attributes)

	if !invoice.Issued.IsZero() {
		root.SetValue("Faktura.Fa.P_1", invoice.Issued.Format("2006-01-02"))
	}
	if invoice.Number != "" {
		root.SetValue("Faktura.Fa.P_2", invoice.Number)
	}
	generationTime := invoice.GenerationTime
	if generationTime.IsZero() {
		generationTime = fg.runTimestamp
	}
	root.SetValue("Faktura.Naglowek.DataWytworzeniaFa", generationTime.Format(time.RFC3339))

	faNode, _ := root.LocateNode("Faktura.Fa")

	if invoice.Correction != nil {
		correctionNode, _ := faNode.CreateChild("DaneFaKorygowanej", false)
		correctionNode.SetValuesFromMap(
			map[string]string{
				"NrFaKorygowanej":       invoice.Correction.RefNo,
				"DataWystFaKorygowanej": invoice.Correction.OriginalIssueDate.Format(time.DateOnly),
				"NrKSeF":                "1",
				"NrKSeFFaKorygowanej":   invoice.Correction.KSeFRefNo,
			},
		)
	}

	for i, item := range invoice.Items {
		faChildNode, _ := faNode.CreateChild("FaWiersz", true)
		rowNo := i + 1
		if item.RowNo > 0 {
			rowNo = item.RowNo
		}
		faChildNode.SetValue("NrWierszaFa", fmt.Sprintf("%d", rowNo))
		if item.Before {
			faChildNode.SetValue("StanPrzed", "1")
		}
		// an empty item may appear on correction invoices when a before / after item
		// does not have any values - thus marking it as empty
		if !item.IsEmpty() {
			faChildNode.SetValue("P_7", item.Description)
			if item.Unit != "" {
				faChildNode.SetValue("P_8A", item.Unit)
			}
			faChildNode.SetValue("P_8B", money.RenderAmountFromCurrencyUnits(item.Quantity.Amount, uint8(item.Quantity.DecimalPlaces)))
			if !item.UnitPrice.IsGross {
				faChildNode.SetValue("P_11", money.RenderAmountFromCurrencyUnits(item.Amount().Net, 2))
			} else {
				faChildNode.SetValue("P_11A", money.RenderAmountFromCurrencyUnits(item.Amount().Gross, 2))
			}
			faChildNode.SetValue("P_12", item.UnitPrice.Vat.Description)
			unitPriceField := "P_9B"
			if !item.UnitPrice.IsGross {
				unitPriceField = "P_9A"
			}
			faChildNode.SetValue(unitPriceField, money.RenderAmountFromCurrencyUnits(item.UnitPrice.Amount, uint8(item.UnitPrice.DecimalPlaces)))
			faChildNode.SetValuesFromMap(item.Attributes)
		}
		fieldToVatRatesMapping.Accumulate(item)
	}

	fieldToVatRatesMapping.Populate(root)
	faNode.SetValue("P_15", money.RenderAmountFromCurrencyUnits(invoice.Total.Gross, 2))
	if err := root.ApplyOrdering(fg.elementOrdering); err != nil {
		return nil, fmt.Errorf("unable to apply schema order: %v", err)
	}

	if invoice.Type != "" {
		faNode.SetValue("RodzajFaktury", invoice.Type)
	}

	fieldToVatRatesMapping.Zero()

	if err := fg.hooks.PostProcess(root); err != nil {
		return nil, err
	}

	return root, nil
}
