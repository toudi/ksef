package income

import (
	"errors"
	"fmt"
	"ksef/internal/invoicesdb/jpk/processors/vat"
	"ksef/internal/invoicesdb/jpk/types"

	"github.com/beevik/etree"
)

const (
	xpathItem = "//Faktura/Fa/FaWiersz"
)

var errUnableToCalculateVAT = errors.New("unable to calculate VAT for invoice item")

func ProcessVatInfo(invoiceXML *etree.Document, invoice *types.Invoice) error {
	items := invoiceXML.FindElements(xpathItem)

	vatCalculator := &vat.Calculator{}

	for _, item := range items {
		vatInfo, err := vatCalculator.GetVat(item)
		if err != nil {
			return errors.Join(errUnableToCalculateVAT, err)
		}

		if _, exists := invoice.VATByRate[vatInfo.VatRate]; !exists {
			invoice.VATByRate[vatInfo.VatRate] = &types.VATInfo{}
		}

		invoice.VATByRate[vatInfo.VatRate].Base = invoice.VATByRate[vatInfo.VatRate].Base.Add(vatInfo.NetAmount)
		invoice.VATByRate[vatInfo.VatRate].Vat = invoice.VATByRate[vatInfo.VatRate].Base.Add(vatInfo.VatAmount)

		invoice.VAT.Base = invoice.VAT.Base.Add(vatInfo.NetAmount)
		invoice.VAT.Vat = invoice.VAT.Vat.Add(vatInfo.VatAmount)
	}

	for vatRate, amount := range invoice.VATByRate {
		switch vatRate {
		case "zw":
			invoice.Attributes["K_10"] = amount.Base.Format(2)
		case "np I":
			invoice.Attributes["K_11"] = amount.Base.Format(2)
		case "np II":
			invoice.Attributes["K_12"] = amount.Base.Format(2)
		case "0":
			invoice.Attributes["K_13"] = amount.Base.Format(2)
		case "5":
			invoice.Attributes["K_15"] = amount.Base.Format(2)
			invoice.Attributes["K_16"] = amount.Vat.Format(2)
		case "7", "8":
			invoice.Attributes["K_17"] = amount.Base.Format(2)
			invoice.Attributes["K_18"] = amount.Vat.Format(2)
		case "22", "23":
			invoice.Attributes["K_19"] = amount.Base.Format(2)
			invoice.Attributes["K_20"] = amount.Vat.Format(2)
		default:
			return fmt.Errorf("unhandled tax rate: %s", vatRate)
		}
	}

	return nil
}
