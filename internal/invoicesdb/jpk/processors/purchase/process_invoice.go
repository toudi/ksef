package purchase

import (
	"ksef/internal/invoicesdb/jpk/interfaces"
	"ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"maps"

	"github.com/beevik/etree"
)

const (
	xpathIssued = "//Faktura/Fa/P_1"
)

func ProcessInvoice(
	dest *types.Invoice,
	invoiceXML *etree.Document,
	registryInvoice *monthlyregistry.Invoice,
	manager interfaces.JPKManager,
) error {
	maps.Copy(dest.Attributes, map[string]string{
		"DataZakupu": invoiceXML.FindElement(xpathIssued).Text(),
	})
	return nil
}
