package fa_3_1

import (
	"ksef/internal/interfaces"
	"ksef/internal/sei/generators/fa"
	"ksef/internal/xml"
	"time"
)

func GeneratorFactory() interfaces.Generator {
	return fa.New(
		fa.WithCommonData(map[string]string{
			"Faktura.#xmlns:xsi":                          "http://www.w3.org/2001/XMLSchema-instance",
			"Faktura.#xmlns:xsd":                          "http://www.w3.org/2001/XMLSchema",
			"Faktura.#xmlns":                              "http://crd.gov.pl/wzor/2025/06/25/13775/",
			"Faktura.Naglowek.KodFormularza":              "FA",
			"Faktura.Naglowek.KodFormularza#kodSystemowy": "FA (3)",
			"Faktura.Naglowek.KodFormularza#wersjaSchemy": "1-0E",
			"Faktura.Naglowek.WariantFormularza":          "3",
			"Faktura.Naglowek.SystemInfo":                 "WSI Pegasus",
			"Faktura.Podmiot2.JST":                        "2",
			"Faktura.Podmiot2.GV":                         "2",
			"Faktura.Podmiot2.Adres.KodKraju":             "PL",
			"Faktura.Fa.Adnotacje.P_16":                   "2", // metoda kasowa
			"Faktura.Fa.Adnotacje.P_17":                   "2", // samofakturowanie
			"Faktura.Fa.Adnotacje.P_18":                   "2", // odwrotne obciążenie
			"Faktura.Fa.Adnotacje.P_18A":                  "2", // mechanizm podzielonej płatności
			"Faktura.Fa.Adnotacje.P_23":                   "2", // trójstronna faktura uproszczona
		}),
		fa.WithElementOrdering(FA_3_1ChildrenOrder),
	)
}

var fa_3_1_hooks = fa.GeneratorHooks{
	PostProcess: func(root *xml.Node) error {
		root.SetValue("Faktura.Naglowek.DataWytworzeniaFa", time.Now().Local().Format(time.RFC3339))
		return nil
	},
}
