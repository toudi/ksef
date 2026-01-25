package jpk

import (
	"ksef/internal/invoicesdb/jpk/interfaces"
	"ksef/internal/invoicesdb/jpk/processors/income"
	"ksef/internal/invoicesdb/jpk/processors/purchase"
	"ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/money"
	"strconv"

	"github.com/beevik/etree"
)

const (
	xpathPreamble = "//Faktura/Fa"
	xpathItems    = xpathPreamble + "/FaWiersz"
)

var incomeInvoiceProcessors = []func(invoice *etree.Document, dest *types.Invoice) error{
	income.ProcessBuyer,
	income.ProcessInvoice,
	income.ProcessVatInfo,
}

var purchaseInvoiceProcessors = []func(
	dest *types.Invoice, xmlInvoice *etree.Document,
	registryInvoice *monthlyregistry.Invoice,
	manager interfaces.JPKManager,
) error{
	purchase.ProcessInvoice,
	purchase.ProcessVatInfo,
}

func (j *JPK) AddIncome(xmlInvoice *etree.Document, invoice *monthlyregistry.Invoice) error {
	// no point in reporting a failed attempt
	if invoice.KSeFRefNo == "" && len(invoice.UploadErrors) > 0 {
		return nil
	}

	income := &types.Invoice{
		Attributes: map[string]string{
			"LpSprzedazy":    strconv.Itoa(len(j.Income) + 1),
			"NrKSeF":         invoice.KSeFRefNo,
			"DowodSprzedazy": invoice.RefNo,
		},
		VAT: &types.VATInfo{
			Base: money.MonetaryValue{},
			Vat:  money.MonetaryValue{},
		},
		VATByRate: make(map[string]*types.VATInfo),
	}

	for _, processor := range incomeInvoiceProcessors {
		if err := processor(xmlInvoice, income); err != nil {
			return err
		}
	}

	j.Income = append(j.Income, income)

	return nil
}

func (j *JPK) AddReceived(xmlInvoice *etree.Document, invoice *monthlyregistry.Invoice) error {
	purchase := &types.Invoice{
		Attributes: map[string]string{
			"LpZakupu":    strconv.Itoa(len(j.Purchase) + 1),
			"NrKSeF":      invoice.KSeFRefNo,
			"DowodZakupu": invoice.RefNo,
		},
		VAT:       &types.VATInfo{},
		VATByRate: make(map[string]*types.VATInfo),
	}

	for _, processor := range purchaseInvoiceProcessors {
		if err := processor(purchase, xmlInvoice, invoice, j.manager); err != nil {
			return err
		}
	}

	j.Purchase = append(j.Purchase, purchase)
	return nil
}
