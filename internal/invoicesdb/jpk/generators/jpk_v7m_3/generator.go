package jpk_v7m_3

import (
	"ksef/internal/xml"
	"time"
)

// if there is more than one generator than maybe we need to extract it to a common interface or something
func Document() *xml.Node {
	root := &xml.Node{Name: "JPK"}
	commonData := map[string]string{
		"JPK.#xmlns":                                                   "http://crd.gov.pl/wzor/2025/12/19/14090/",
		"JPK.#xmlns:xsi":                                               "http://www.w3.org/2001/XMLSchema-instance",
		"JPK.Naglowek.KodFormularza":                                   "JPK_VAT",
		"JPK.Naglowek.KodFormularza#kodSystemowy":                      "JPK_V7M (3)",
		"JPK.Naglowek.KodFormularza#wersjaSchemy":                      "1-0E",
		"JPK.Naglowek.WariantFormularza":                               "3",
		"JPK.Naglowek.NazwaSystemu":                                    "WSI Pegasus",
		"JPK.Naglowek.CelZlozenia#poz":                                 "P_7",
		"JPK.Naglowek.CelZlozenia":                                     "1",
		"JPK.Deklaracja.Naglowek.KodFormularzaDekl":                    "VAT-7",
		"JPK.Deklaracja.Naglowek.KodFormularzaDekl#kodSystemowy":       "VAT-7 (23)",
		"JPK.Deklaracja.Naglowek.KodFormularzaDekl#kodPodatku":         "VAT",
		"JPK.Deklaracja.Naglowek.KodFormularzaDekl#rodzajZobowiazania": "Z",
		"JPK.Deklaracja.Naglowek.KodFormularzaDekl#wersjaSchemy":       "1-0E",
		"JPK.Deklaracja.Naglowek.WariantFormularzaDekl":                "23",
		"JPK.Deklaracja.Pouczenia":                                     "1",
	}

	root.SetValuesFromMap(commonData)

	root.SetValue("JPK.Naglowek.DataWytworzeniaJPK", time.Now().Format(time.RFC3339))

	return root
}
