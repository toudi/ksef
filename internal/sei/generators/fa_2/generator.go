package fa_2

import (
	"ksef/internal/interfaces"
	"ksef/internal/sei/generators/fa"
)

func GeneratorFactory() interfaces.Generator {
	return fa.New(
		fa.WithCommonData(map[string]string{
			"Faktura.#xmlns:xsi":                          "http://www.w3.org/2001/XMLSchema-instance",
			"Faktura.#xmlns:xsd":                          "http://www.w3.org/2001/XMLSchema",
			"Faktura.#xmlns":                              "http://crd.gov.pl/wzor/2023/06/29/12648/",
			"Faktura.Naglowek.KodFormularza":              "FA",
			"Faktura.Naglowek.KodFormularza#kodSystemowy": "FA (2)",
			"Faktura.Naglowek.KodFormularza#wersjaSchemy": "1-0E",
			"Faktura.Naglowek.WariantFormularza":          "2",
			"Faktura.Naglowek.SystemInfo":                 "WSI Pegasus",
		}),
		fa.WithElementOrdering(FA_2ChildrenOrder),
	)
}
