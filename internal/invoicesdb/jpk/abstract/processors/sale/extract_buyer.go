package sale

import (
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

const (
	recipientPrefix = "//Faktura/Podmiot2/"
	recipientNIP    = recipientPrefix + "DaneIdentyfikacyjne/NIP"
	recipientName   = recipientPrefix + "DaneIdentyfikacyjne/Nazwa"
)

func ExtractBuyer(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	salesRow *types.SaleItem,
) error {
	recipientID := "BRAK"
	recipientNIPNode := doc.FindElement(recipientNIP)

	if recipientNIPNode != nil {
		recipientID = recipientNIPNode.Text()
	}

	salesRow.Buyer.NIP = recipientID
	salesRow.Buyer.Name = doc.FindElement(recipientName).Text()

	return nil
}
