package jpk_v7m_3

import (
	"errors"
	"ksef/internal/invoicesdb/jpk/abstract"
	"ksef/internal/invoicesdb/jpk/generators/interfaces"
	"ksef/internal/invoicesdb/jpk/manager"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/money"
	"ksef/internal/xml"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
)

const (
	subjectTypeIndividual           = "f"
	declarationDetailedInfoBasePath = "JPK.Deklaracja.PozycjeSzczegolowe."
	EdtNamespace                    = "http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2022/09/13/eD/DefinicjeTypy/"
)

var individualEDTFields = []string{"NIP", "ImiePierwsze", "Nazwisko", "DataUrodzenia"}

var commonData = map[string]string{
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

type jpk_v7m_3_generator struct {
	report  *abstract.MonthlyReport
	month   time.Time
	manager *manager.JPKManager
}

func New(manager *manager.JPKManager, month time.Time) interfaces.JPKGenerator {
	return &jpk_v7m_3_generator{
		report:  abstract.NewMonthlyReport(manager),
		month:   month,
		manager: manager,
	}
}

func (g *jpk_v7m_3_generator) Document() (*xml.Node, error) {
	root := &xml.Node{Name: "JPK"}

	root.SetValuesFromMap(commonData)
	root.SetValue("JPK.Naglowek.DataWytworzeniaJPK", time.Now().Format(time.RFC3339))

	sjs := g.manager.GetSettings()
	if sjs != nil {
		if sjs.FormMeta.IRSCode > 0 {
			root.SetValue("JPK.Naglowek.KodUrzedu", strconv.Itoa(sjs.FormMeta.IRSCode))
		}
		if sjs.FormMeta.SystemName != "" {
			root.SetValue("JPK.Naglowek.NazwaSystemu", sjs.FormMeta.SystemName)
		}
	} else {
		return nil, errors.New("JPK settings not defined")
	}

	root.SetValuesFromMap(
		map[string]string{
			"JPK.Naglowek.Rok":     g.month.Format("2006"),
			"JPK.Naglowek.Miesiac": g.month.Format("01"),
			"JPK.Podmiot1#rola":    "Podatnik",
		},
	)

	// some of these fields will be overwritten by row processing functions, but we
	// do have to initialize them to zero values - they are mandatory.
	for node_name, default_value := range JPK_V7M_3RequiredDefaults {
		// let's extract the prefix to check if it is contained in the array
		// nodes. if so - we cannot apply it here.
		node_name_parts := strings.Split(node_name, ".")
		node_prefix := strings.Join(node_name_parts[:len(node_name_parts)-1], ".")
		if JPK_V7M_3ArrayElements[node_prefix] {
			continue
		}
		root.SetValue(node_name, default_value)
	}

	// now follow up with defaults from subject settings
	for node_name, default_value := range sjs.FormMeta.Defaults {
		root.SetValue(node_name, default_value)
	}

	if sjs.FormMeta.Subject != nil {
		isIndividual := strings.ToLower(sjs.FormMeta.Subject.SubjectType) == subjectTypeIndividual
		subjectTypeName := "OsobaNiefizyczna"
		if isIndividual {
			root.SetValue("JPK.#xmlns:edt", EdtNamespace)
			subjectTypeName = "OsobaFizyczna"
		}
		for subjectField, fieldValue := range sjs.FormMeta.Subject.Data {
			if isIndividual && slices.Contains(individualEDTFields, subjectField) {
				subjectField = "edt:" + subjectField
			}
			root.SetValue("JPK.Podmiot1."+subjectTypeName+"."+subjectField, fieldValue)
		}
	}

	// field to amount is an aggregation helper that will allow us to dynamically add some amounts
	fieldToAmount := &fieldToAmountRegistry{}

	ewidencja, _ := root.CreateChild("Ewidencja", false)

	// populate SprzedazWiersz ..
	for saleRowIdx, saleRow := range g.report.Sales.Rows {
		sprzedazWiersz, _ := ewidencja.CreateChild("SprzedazWiersz", true)
		sprzedazWiersz.SetValue("LpSprzedazy", strconv.Itoa(saleRowIdx+1))

		if err := populateSprzedazWiersz(sprzedazWiersz, saleRow, fieldToAmount); err != nil {
			return nil, err
		}
	}
	// .. and SprzedazCtrl
	sprzedazWierszCtrl, _ := ewidencja.CreateChild("SprzedazCtrl", false)
	sprzedazWierszCtrl.SetValuesFromMap(map[string]string{
		"LiczbaWierszySprzedazy": strconv.Itoa(len(g.report.Sales.Rows)),
		"PodatekNalezny":         g.report.Sales.VATAmounts.Total.Vat.Format(2),
	})

	// populate ZakupWiersz ..
	for purchaseIdx, purchaseRow := range g.report.Purchase.Rows {
		zakupWiersz, _ := ewidencja.CreateChild("ZakupWiersz", true)
		zakupWiersz.SetValue("LpZakupu", strconv.Itoa(purchaseIdx+1))

		if err := populateZakupWiersz(zakupWiersz, purchaseRow, fieldToAmount); err != nil {
			return nil, err
		}
	}
	// .. and ZakupCtrl
	zakupWierszCtrl, _ := ewidencja.CreateChild("ZakupCtrl", false)
	zakupWierszCtrl.SetValuesFromMap(map[string]string{
		"LiczbaWierszyZakupu": strconv.Itoa(len(g.report.Purchase.Rows)),
		"PodatekNaliczony":    g.report.Purchase.VATAmounts.Total.Vat.Format(2),
	})

	// now we can populate accumulated values:
	for field, amount := range *fieldToAmount {
		root.SetValue(declarationDetailedInfoBasePath+field, amount.Format(0))
	}
	// whew that was a lot of work.
	// now we can populate PozycjeSzczeglowe fields of Deklaracja section.
	// Calculate P_37 (total base)
	// according to schema:
	// P_10, P_11, P_13, P_15, P_17, P_19, P_21, P_22, P_23, P_25, P_27, P_29, P_31
	accumulated := fieldToAmount.accumulate(
		[]string{
			"P_10", "P_11", "P_13", "P_15", "P_17", "P_19", "P_21", "P_22", "P_23", "P_25",
			"P_27", "P_29", "P_31",
		},
	)
	root.SetValue(declarationDetailedInfoBasePath+"P_37", accumulated.Format(0))

	// Calculate P_38 (total output VAT)
	// according to schema:
	// P_16, P_18, P_20, P_24, P_26, P_28, P_30, P_32, P_33, P_34 pomniejszona o kwotę z P_35, P_36 i P_360
	vatSum := fieldToAmount.accumulate(
		[]string{"P_16", "P_18", "P_20", "P_24", "P_26", "P_28", "P_30", "P_32", "P_33", "P_34"},
	)
	subtractSum := fieldToAmount.accumulate(
		[]string{"P_35", "P_36", "P_360"},
	)
	vatSum = vatSum.Add(money.MonetaryValue{
		Amount:        subtractSum.Amount * -1,
		DecimalPlaces: subtractSum.DecimalPlaces,
	})
	root.SetValue(declarationDetailedInfoBasePath+"P_38", vatSum.Format(0))

	// Calculate P_48
	// according to schema:
	// P_39, P_41, P_43, P_44, P_45, P_46 i P_47
	vatSum = fieldToAmount.accumulate(
		[]string{"P_39", "P_41", "P_43", "P_44", "P_45", "P_46", "P_47"},
	)
	root.SetValue(declarationDetailedInfoBasePath+"P_48", vatSum.Format(0))

	if err := root.ApplyOrdering(JPK_V7M_3ChildrenOrder); err != nil {
		return nil, err
	}

	return root, nil
}

func (g *jpk_v7m_3_generator) ProcessInvoice(
	invoice *monthlyregistry.Invoice,
	document *etree.Document,
) error {
	if invoice.Type == monthlyregistry.InvoiceTypeIssued {
		return g.report.Sales.ProcessInvoice(invoice, document)
	}
	return g.report.Purchase.ProcessInvoice(invoice, document)
}
