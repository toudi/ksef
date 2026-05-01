package jpk_v7m_3

import (
	absTypes "ksef/internal/invoicesdb/jpk/abstract/types"
	"ksef/internal/xml"
	"strings"
)

const zakupWierszPrefix = "JPK.Ewidencja.ZakupWiersz."

func populateZakupWiersz(dest *xml.Node, row *absTypes.PurchaseItem, declaration *fieldToAmountRegistry) error {
	data := map[string]string{
		"DataZakupu":    row.Date,
		"NrDostawcy":    row.Seller.NIP,
		"NazwaDostawcy": row.Seller.Name,
		"DowodZakupu":   row.RefNo,
		"NrKSeF":        row.KSeFRefNo,
	}

	// populate required defaults
	for field, _default := range JPK_V7M_3RequiredDefaults {
		if strings.HasPrefix(field, zakupWierszPrefix) {
			fieldName := strings.TrimPrefix(field, zakupWierszPrefix)
			data[fieldName] = _default
		}
	}

	for purchaseAttributes, amounts := range row.VATAmounts.ByAttributes {
		formFields := FormFields{BaseField: "K_42", VatField: "K_43"}
		if purchaseAttributes.FixedAssets {
			formFields.BaseField = "K_40"
			formFields.VatField = "K_41"
		}
		data[formFields.BaseField] = amounts.Base.Format(2)
		data[formFields.VatField] = amounts.Vat.Format(2)

		declaration.Add(formFields, amounts)
	}

	dest.SetValuesFromMap(data)
	return nil
}
