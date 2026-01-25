package jpk

import (
	"fmt"
	"ksef/internal/invoicesdb/jpk/generators/jpk_v7m_3"
	"ksef/internal/xml"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	JPKPurchase = "JPK.Ewidencja.ZakupWiersz."
	JPKIncome   = "JPK.Ewidencja.SprzedazWiersz."
)

func (j *JPK) Save(output string) error {
	root := jpk_v7m_3.Document()
	j.populatePreamble(root)
	// now we can iterate over income and purchase invoices and generate appropriate keys
	invoiceRows, _ := root.CreateChild("Ewidencja", false)

	for _, income := range j.Income {
		incomeRow, _ := invoiceRows.CreateChild("SprzedazWiersz", true)
		// copy over the defaults
		for field, value := range jpk_v7m_3.JPK_V7M_3RequiredDefaults {
			if strings.HasPrefix(field, JPKIncome) {
				incomeRow.SetValue(strings.TrimPrefix(field, JPKIncome), value)
			}
		}
		incomeRow.SetValuesFromMap(income.Attributes)
		j.IncomeCtrl.VAT = j.IncomeCtrl.VAT.Add(income.VAT.Vat)
	}
	// generate SprzedazCtrl row that contains aggregates and counter
	incomeCtrl, _ := root.CreateChild("Ewidencja.SprzedazCtrl", false)

	incomeCtrl.SetValuesFromMap(map[string]string{
		"LiczbaWierszySprzedazy": strconv.Itoa(len(j.Income)),
		"PodatekNalezny":         j.IncomeCtrl.VAT.Format(2),
	})

	for _, purchase := range j.Purchase {
		purchaseRow, _ := invoiceRows.CreateChild("ZakupWiersz", true)
		// copy over the defaults
		for field, value := range jpk_v7m_3.JPK_V7M_3RequiredDefaults {
			if strings.HasPrefix(field, JPKPurchase) {
				purchaseRow.SetValue(strings.TrimPrefix(field, JPKPurchase), value)
			}
		}
		purchaseRow.SetValuesFromMap(purchase.Attributes)
		j.PurchaseCtrl.VAT = j.PurchaseCtrl.VAT.Add(purchase.VAT.Vat)
	}
	// generate SprzedazCtrl row that contains aggregates and counter
	purchaseCtrl, _ := root.CreateChild("Ewidencja.ZakupCtrl", false)

	purchaseCtrl.SetValuesFromMap(map[string]string{
		"LiczbaWierszyZakupow": strconv.Itoa(len(j.Purchase)),
		"PodatekNaliczony":     j.PurchaseCtrl.VAT.Format(2),
	})

	return j.writeToFile(root, output)
}

func (j *JPK) populatePreamble(root *xml.Node) {
	if j.sjs.FormMeta.IRSCode > 0 {
		root.SetValue("JPK.Naglowek.KodUrzedu", strconv.Itoa(j.sjs.FormMeta.IRSCode))
	}
	if j.sjs.FormMeta.SystemName != "" {
		root.SetValue("JPK.Naglowek.NazwaSystemu", j.sjs.FormMeta.SystemName)
	}
	// data/<env>/<nip><year>/<month>
	//                   -2    -1
	pathParts := strings.Split(j.path, string(filepath.Separator))
	root.SetValuesFromMap(
		map[string]string{
			"JPK.Naglowek.Rok":     pathParts[len(pathParts)-2],
			"JPK.Naglowek.Miesiac": pathParts[len(pathParts)-1],
			"JPK.Podmiot1#rola":    "Podatnik",
		},
	)

	// some of these fields will be overwritten by row processing functions, but we
	// do have to initialize them to zero values - they are mandatory.
	for node_name, default_value := range jpk_v7m_3.JPK_V7M_3RequiredDefaults {
		// let's extract the prefix to check if it is contained in the array
		// nodes. if so - we cannot apply it here.
		node_name_parts := strings.Split(node_name, ".")
		node_prefix := strings.Join(node_name_parts[:len(node_name_parts)-1], ".")
		if jpk_v7m_3.JPK_V7M_3ArrayElements[node_prefix] {
			continue
		}
		root.SetValue(node_name, default_value)
	}

	// now follow up with defaults from subject settings
	for node_name, default_value := range j.sjs.FormMeta.Defaults {
		root.SetValue(node_name, default_value)
	}

	if j.sjs.FormMeta.Subject != nil {
		for subjectType, typeValues := range j.sjs.FormMeta.Subject {
			for keyName, keyValue := range typeValues.(map[string]any) {
				root.SetValue("JPK.Podmiot1."+subjectType+"."+keyName, fmt.Sprintf("%v", keyValue))
			}
		}
	}
}

func (j *JPK) writeToFile(root *xml.Node, outputDir string) error {
	var err error
	if err = os.MkdirAll(outputDir, 0775); err != nil {
		return err
	}
	outputFilename := "jpk-v7m.xml"
	// TODO: properly recognize if we're creating a correction. if so,
	// set the following path:
	// root.SetValue("JPK.Naglowek.CelZlozenia", "2")
	if err = root.ApplyOrdering(jpk_v7m_3.JPK_V7M_3ChildrenOrder); err != nil {
		return err
	}
	writer, err := os.Create(filepath.Join(outputDir, outputFilename))
	if err != nil {
		return err
	}
	defer writer.Close()
	return root.DumpToWriter(writer, 0)
}
