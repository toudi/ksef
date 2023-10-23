package fa_1

import (
	"fmt"
	"ksef/common"
	"ksef/common/xml"
	"strings"
)

type FA1Generator struct {
	commonData          map[string]string
	invoices            []*xml.Node
	currentInvoiceIndex int
}

func (fg *FA1Generator) createNewInvoice(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == sectionInvoice
}
func (fg *FA1Generator) isCommonData(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == "faktura" || sectionLower == "faktura.naglowek" || strings.HasPrefix(sectionLower, "faktura.podmiot1")
}

func (fg *FA1Generator) LineHandler(section string, data map[string]string) error {
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

func (fg *FA1Generator) newInvoice() *xml.Node {
	var root = &xml.Node{Name: "Faktura"}

	// root.SetValue("Faktura.#xmlns:etd", "http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2021/06/09/eD/DefinicjeTypy/")

	for key, value := range fg.commonData {
		root.SetValue(key, value)
	}

	return root
}

func (fg *FA1Generator) Save(dest string) error {
	var err error
	var i int
	var invoice *xml.Node

	for i, invoice = range fg.invoices {
		if err = FA_1(invoice, fmt.Sprintf("%s/faktura_%d.xml", dest, i)); err != nil {
			return fmt.Errorf("unable to generate invoice %d: %v", i, err)
		}
	}
	return nil
}

func (fg *FA1Generator) IssuerTIN() string {
	return fg.commonData["Faktura.Podmiot1.DaneIdentyfikacyjne.NIP"]
}

func (fg *FA1Generator) InvoiceToXMLTree(invoice *common.Invoice) (*xml.Node, error) {
	return nil, fmt.Errorf("not implemented")
}

func GeneratorFactory() common.Generator {
	return &FA1Generator{
		currentInvoiceIndex: -1,
		commonData: map[string]string{
			"Faktura.#xmlns:xsi":                          "http://www.w3.org/2001/XMLSchema-instance",
			"Faktura.#xmlns:xsd":                          "http://www.w3.org/2001/XMLSchema",
			"Faktura.#xmlns":                              "http://crd.gov.pl/wzor/2021/11/29/11089/",
			"Faktura.Naglowek.KodFormularza":              "FA",
			"Faktura.Naglowek.KodFormularza#kodSystemowy": "FA (1)",
			"Faktura.Naglowek.KodFormularza#wersjaSchemy": "1-0E",
			"Faktura.Naglowek.WariantFormularza":          "1",
			"Faktura.Naglowek.SystemInfo":                 "WSI Pegasus",
		},
	}
}
