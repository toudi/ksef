package jpk_v7m_3

import (
	"errors"
	absTypes "ksef/internal/invoicesdb/jpk/abstract/types"
	"ksef/internal/invoicesdb/jpk/types"
	"ksef/internal/money"
	"ksef/internal/xml"
	"strings"
)

const sprzedazWierszPrefix = "JPK.Ewidencja.SprzedazWiersz."

// mapping of sale attributes to their schema-specific declaration fields.
var saleAttributesToDeclarationFields = map[types.SaleAttributes]FormFields{
	{VatRate: types.VatRateZw}:      {BaseField: "K_10"},
	{VatRate: types.VatRateNpI}:     {BaseField: "K_11"},
	{VatRate: types.VatRateNpII}:    {BaseField: "K_12"},
	{VatRate: types.VatRateZeroKR}:  {BaseField: "K_13"},
	{VatRate: types.VatRateZeroWDT}: {BaseField: "K_21"},
	{VatRate: types.VatRateZeroEX}:  {BaseField: "K_22"},
	{VatRate: types.VatRate5}:       {BaseField: "K_15", VatField: "K_16"},
	{VatRate: types.VatRate7}:       {BaseField: "K_17", VatField: "K_18"},
	{VatRate: types.VatRate8}:       {BaseField: "K_17", VatField: "K_18"},
	{VatRate: types.VatRate22}:      {BaseField: "K_19", VatField: "K_20"},
	{VatRate: types.VatRate23}:      {BaseField: "K_19", VatField: "K_20"},
}

func populateSprzedazWiersz(dest *xml.Node, row *absTypes.SaleItem, declaration *fieldToAmountRegistry) error {
	data := map[string]string{
		"DataWystawienia":  row.IssueDate,
		"DataSprzedazy":    row.SaleDate,
		"NrKontrahenta":    row.Buyer.NIP,
		"NazwaKontrahenta": row.Buyer.Name,
	}

	for field, defaultValue := range JPK_V7M_3RequiredDefaults {
		if strings.HasPrefix(field, sprzedazWierszPrefix) {
			fieldName := strings.TrimPrefix(field, sprzedazWierszPrefix)
			data[fieldName] = defaultValue
		}
	}

	// now that we've got required values populated let's populate the ones that we
	// understand.
	accumulator := &amountsAccmulator{}

	for vatRate, amounts := range row.VATAmounts.ByRate {
		fields, exists := saleAttributesToDeclarationFields[types.SaleAttributes{VatRate: vatRate}]
		if !exists {
			return errors.New("unexpected vat rate")
		}
		accumulator.Add(fields, amounts)
		declaration.Add(fields, amounts)
	}

	for fieldName, amount := range accumulator.values {
		if amount.Amount > 0 {
			data[fieldName] = amount.Format(2)
		}
	}

	dest.SetValuesFromMap(data)
	return nil
}

type amountsAccmulator struct {
	values map[string]money.MonetaryValue
}

func (aa *amountsAccmulator) Add(fields FormFields, vatInfo *types.VATInfo) {
	if _, exists := aa.values[fields.BaseField]; !exists {
		aa.values[fields.BaseField] = money.MonetaryValue{}
	}
	aa.values[fields.BaseField] = aa.values[fields.BaseField].Add(vatInfo.Base)
	if fields.VatField != "" {
		if _, exists := aa.values[fields.VatField]; !exists {
			aa.values[fields.VatField] = money.MonetaryValue{}
		}
		aa.values[fields.VatField] = aa.values[fields.VatField].Add(vatInfo.Vat)
	}
}
