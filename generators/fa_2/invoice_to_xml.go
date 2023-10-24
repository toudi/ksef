package fa_2

import (
	"fmt"
	"ksef/common"
	"ksef/common/xml"
	"time"
)

// this function is meant to be used when you have a invoice.Invoice object as your source
// and want to output the xml.Node object as your output, which is then used to generate
// metadata and so on.
// if this is not what you wish to do, then please use the line-based generator approach (either CSV or XLSX)
func (fg *FA2Generator) InvoiceToXMLTree(invoice *common.Invoice) (*xml.Node, error) {
	var root = &xml.Node{Name: "Faktura"}

	root.SetValuesFromMap(invoice.Attributes)

	if !invoice.Issued.IsZero() {
		root.SetValue("Faktura.Fa.P_1", invoice.Issued.Format("2006-01-02"))
	}
	if invoice.Number != "" {
		root.SetValue("Faktura.Fa.P_2", invoice.Number)
	}
	root.SetValue("Faktura.Naglowek.DataWytworzeniaFa", fg.runTimestamp.Format(time.RFC3339))

	faNode, _ := root.LocateNode("Faktura.Fa")

	for i, item := range invoice.Items {
		faChildNode, _ := faNode.CreateChild("FaWiersz", true)
		faChildNode.SetValue("NrWierszaFa", fmt.Sprintf("%d", i+1))
		faChildNode.SetValue("P_7", item.Description)
		if item.Unit != "" {
			faChildNode.SetValue("P_8A", item.Unit)
		}
		faChildNode.SetValue("P_8B", common.RenderFloatNumber(item.Quantity))
		faChildNode.SetValue("P_11", common.RenderAmountFromCurrencyUnits(item.Amount().Net, 2))
		faChildNode.SetValue("P_12", item.UnitPrice.Vat.Description)
		if !item.UnitPrice.IsGross {
			faChildNode.SetValue("P_9A", common.RenderAmountFromCurrencyUnits(item.UnitPrice.Value, 2))
		} else {
			faChildNode.SetValue("P_9B", common.RenderAmountFromCurrencyUnits(item.UnitPrice.Value, 2))
		}
		faChildNode.SetValuesFromMap(item.Attributes)
		fieldToVatRatesMapping.Accumulate(item)
	}

	fieldToVatRatesMapping.Populate(root)
	faNode.SetValue("P_15", common.RenderAmountFromCurrencyUnits(invoice.Total.Gross, 2))
	if err := root.ApplyOrdering(FA_2ChildrenOrder); err != nil {
		return nil, fmt.Errorf("unable to apply schema order: %v", err)
	}

	fieldToVatRatesMapping.Zero()

	return root, nil
}
