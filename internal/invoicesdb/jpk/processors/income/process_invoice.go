package income

import (
	"ksef/internal/invoicesdb/jpk/types"
	"maps"

	"github.com/beevik/etree"
)

const (
	xpathIssued = "//Faktura/Fa/P_1"
	xpathSale   = "//Faktura/Fa/P_6"
)

func ProcessInvoice(invoiceXML *etree.Document, invoice *types.Invoice) error {
	maps.Copy(invoice.Attributes, map[string]string{
		"DataWystawienia": invoiceXML.FindElement(xpathIssued).Text(),
		"DataSprzedazy":   invoiceXML.FindElement(xpathSale).Text(),
	})
	return nil
}
