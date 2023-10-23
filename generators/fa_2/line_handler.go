package fa_2

import "ksef/common/xml"

func (fg *FA2Generator) LineHandler(section string, data map[string]string) error {
	if fg.isCommonData(section) {
		if fg.commonData == nil {
			fg.commonData = make(map[string]string)
		}
		for key, value := range data {
			fg.commonData[section+"."+key] = value
		}
		return nil
	}
	if fg.createNewInvoice(section) {
		fg.invoices = append(fg.invoices, fg.newInvoice())
		fg.currentInvoiceIndex += 1
	}
	var node *xml.Node = fg.invoices[fg.currentInvoiceIndex]
	node.SetData(section, data)
	return nil
}
