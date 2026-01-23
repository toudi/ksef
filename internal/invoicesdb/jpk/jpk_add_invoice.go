package jpk

import (
	"fmt"
	"ksef/internal/invoicesdb/jpk/generators/jpk_v7m_3"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (j *JPK) AddIncome(invoice *monthlyregistry.XMLInvoice) error {
	// fmt.Printf("invoice: %+v\n", invoice)
	return nil
}

func (j *JPK) AddReceived(invoice *monthlyregistry.XMLInvoice) error {
	return nil
}

func (j *JPK) Save(output string) error {
	var err error
	if err = os.MkdirAll(filepath.Dir(output), 0775); err != nil {
		return err
	}
	root := jpk_v7m_3.Document()
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
	root.SetValuesFromMap(jpk_v7m_3.JPK_V7M_3RequiredDefaults)
	if j.sjs.FormMeta.Subject != nil {
		for subjectType, typeValues := range j.sjs.FormMeta.Subject {
			for keyName, keyValue := range typeValues.(map[string]any) {
				root.SetValue("JPK.Podmiot1."+subjectType+"."+keyName, fmt.Sprintf("%v", keyValue))
			}
		}
	}
	if err = root.ApplyOrdering(jpk_v7m_3.JPK_V7M_3ChildrenOrder); err != nil {
		return err
	}
	writer, err := os.Create(output)
	if err != nil {
		return err
	}
	defer writer.Close()
	return root.DumpToWriter(writer, 0)
}
