package fa_2

import (
	"fmt"
	"ksef/common"
	"ksef/common/xml"
)

// this function is meant to be used when you have a invoice.Invoice object as your source
// and want to output the xml.Node object as your output, which is then used to generate
// metadata and so on.
// if this is not what you wish to do, then please use the line-based generator approach (either CSV or XLSX)
func (fg *FA2Generator) InvoiceToXMLTree(invoice *common.Invoice) (*xml.Node, error) {
	var root = &xml.Node{Name: "Faktura"}

	root.SetValue("Faktura.Fa.P_1", invoice.Issued.Format("2006-01-02"))
	root.SetValue("Faktura.Fa.P_2", invoice.Number)

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
		if err := faChildNode.SetData("", item.Attributes); err != nil {
			return nil, fmt.Errorf("unable to set attributes for Faktura.Fa.Fawiersz: %v", err)
		}
		fieldToVatRatesMapping.Accumulate(item)
	}

	fieldToVatRatesMapping.Populate(root)

	return root, nil
}
