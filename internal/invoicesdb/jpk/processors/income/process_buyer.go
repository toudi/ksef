package income

import (
	"ksef/internal/invoicesdb/jpk/types"

	"github.com/beevik/etree"
)

const (
	recipientPrefix = "//Faktura/Podmiot2/"
	recipientNIP    = recipientPrefix + "DaneIdentyfikacyjne/NIP"
	recipientName   = recipientPrefix + "DaneIdentyfikacyjne/Nazwa"
)

func ProcessBuyer(invoiceXML *etree.Document, invoice *types.Invoice) error {
	recipientID := "BRAK"
	recipientNIPNode := invoiceXML.FindElement(recipientNIP)

	if recipientNIPNode != nil {
		recipientID = recipientNIPNode.Text()
	}

	invoice.Attributes["NrKontrahenta"] = recipientID
	invoice.Attributes["NazwaKontrahenta"] = invoiceXML.FindElement(recipientName).Text()

	return nil
}
