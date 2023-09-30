package generators

import (
	"fmt"
	"ksef/common/xml"
	"strings"
)

type FA2Generator struct {
	commonData          map[string]string
	invoices            []*xml.Node
	currentInvoiceIndex int
}

func (fg *FA2Generator) createNewInvoice(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == sectionInvoice
}
func (fg *FA2Generator) isCommonData(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == "faktura" || sectionLower == "faktura.naglowek" || strings.HasPrefix(sectionLower, "faktura.podmiot1")
}

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

func (fg *FA2Generator) newInvoice() *xml.Node {
	var root = &xml.Node{Name: "Faktura"}

	root.SetValue("Faktura.#xmlns:etd", "http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2022/01/05/eD/DefinicjeTypy/")
	root.SetValue("Faktura.#xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	root.SetValue("Faktura.#xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	root.SetValue("Faktura.#xmlns", "http://crd.gov.pl/wzor/2023/06/29/12648/")
	root.SetValue("Faktura.Naglowek.KodFormularza", "FA")
	root.SetValue("Faktura.Naglowek.KodFormularza#kodSystemowy", "FA (2)")
	root.SetValue("Faktura.Naglowek.KodFormularza#wersjaSchemy", "1-0E")
	root.SetValue("Faktura.Naglowek.WariantFormularza", "2")
	root.SetValue("Faktura.Naglowek.SystemInfo", "WSI Pegasus")

	for key, value := range fg.commonData {
		root.SetValue(key, value)
	}

	return root
}

func (fg *FA2Generator) Save(dest string) error {
	var err error
	var i int
	var invoice *xml.Node

	for i, invoice = range fg.invoices {
		if err = FA_2(invoice, fmt.Sprintf("%s/faktura_%d.xml", dest, i)); err != nil {
			return fmt.Errorf("unable to generate invoice %d: %v", i, err)
		}
	}
	return nil
}

func (fg *FA2Generator) Issuer() string {
	return fg.commonData["Faktura.Podmiot1.DaneIdentyfikacyjne.NIP"]
}

func init() {
	var _generator = &FA2Generator{currentInvoiceIndex: -1}
	registerGenerator("fa_2", _generator)
}
