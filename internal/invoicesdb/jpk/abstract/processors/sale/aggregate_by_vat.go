package sale

import (
	"errors"
	"ksef/internal/invoicesdb/jpk/abstract/processors/vat"
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

const (
	xpathItem = "//Faktura/Fa/FaWiersz"
)

var (
	errUnableToCalculateVAT = errors.New("unable to calculate VAT for invoice item")
	errUnrecognizedVatRate  = errors.New("unrecognizable VAT rate")
)

func AggregateAmountsByVATRate(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	salesRow *types.SaleItem,
) error {
	items := doc.FindElements(xpathItem)

	vatCalculator := &vat.Calculator{}

	for _, item := range items {
		vatInfo, err := vatCalculator.GetVat(item)
		if err != nil {
			return errors.Join(errUnableToCalculateVAT, err)
		}

		if err := salesRow.AddAmount(vatInfo); err != nil {
			return err
		}
	}

	return nil
}
