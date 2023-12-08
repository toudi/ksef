package fa_2

import (
	"ksef/internal/interfaces"
	"ksef/internal/sei/constants"
	"ksef/internal/xml"
	"strings"
	"time"
)

type FA2Generator struct {
	commonData   map[string]string
	runTimestamp time.Time
	// whether all the prices are based on net amount
}

func (fg *FA2Generator) createNewInvoice(section string) bool {
	var sectionLower = strings.ToLower(section)

	return sectionLower == constants.SectionInvoice
}
func (fg *FA2Generator) isCommonData(section string) bool {
	var sectionLower = strings.ToLower(section)

	return (sectionLower == constants.SectionInvoiceRoot ||
		sectionLower == constants.SectionInvoiceHeader ||
		sectionLower == constants.SectionInvoiceAnnotations ||
		strings.HasPrefix(sectionLower, constants.SectionInvoiceIssuer))
}

func (fg *FA2Generator) isItemSection(section string) bool {
	return strings.ToLower(section) == constants.SectionInvoiceItemRow
}

func (fg *FA2Generator) newInvoice() *xml.Node {
	var root = &xml.Node{Name: "Faktura"}

	for key, value := range fg.commonData {
		root.SetValue(key, value)
	}

	return root
}

func (fg *FA2Generator) IssuerTIN() string {
	return fg.commonData["Faktura.Podmiot1.DaneIdentyfikacyjne.NIP"]
}

func GeneratorFactory() interfaces.Generator {
	return &FA2Generator{
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
		runTimestamp: time.Now(),
	}
}
