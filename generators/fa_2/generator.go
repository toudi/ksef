package fa_2

import (
	"fmt"
	"ksef/common"
	"ksef/common/xml"
	"strings"
)

type FA2Generator struct {
	commonData          map[string]string
	invoices            []*xml.Node
	currentInvoiceIndex int
	netBased            bool
	state               int
	// whether all the prices are based on net amount
}

func (fg *FA2Generator) createNewInvoice(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == common.SectionInvoice
}
func (fg *FA2Generator) isCommonData(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == "faktura" || sectionLower == "faktura.naglowek" || strings.HasPrefix(sectionLower, "faktura.podmiot")
}

func (fg *FA2Generator) newInvoice() *xml.Node {
	var root = &xml.Node{Name: "Faktura"}

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

func (fg *FA2Generator) IssuerTIN() string {
	return fg.commonData["Faktura.Podmiot1.DaneIdentyfikacyjne.NIP"]
}

func GeneratorFactory() common.Generator {
	return &FA2Generator{
		currentInvoiceIndex: -1,
		commonData: map[string]string{
			"Faktura.#xmlns:xsi":                          "http://www.w3.org/2001/XMLSchema-instance",
			"Faktura.#xmlns:xsd":                          "http://www.w3.org/2001/XMLSchema",
			"Faktura.#xmlns":                              "http://crd.gov.pl/wzor/2023/06/29/12648/",
			"Faktura.Naglowek.KodFormularza":              "FA",
			"Faktura.Naglowek.KodFormularza#kodSystemowy": "FA (2)",
			"Faktura.Naglowek.KodFormularza#wersjaSchemy": "1-0E",
			"Faktura.Naglowek.WariantFormularza":          "2",
			"Faktura.Naglowek.SystemInfo":                 "WSI Pegasus",
		},
	}
}
