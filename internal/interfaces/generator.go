package interfaces

import (
	"ksef/internal/invoice"
	"ksef/internal/xml"
)

type Generator interface {
	InvoiceToXMLTree(invoice *invoice.Invoice) (*xml.Node, error)
	LineHandler(invoice *invoice.Invoice, section string, data map[string]string, invoiceReady func() error) error
	// method that returns issuer tax identification number from the invoice object
	IssuerTIN() string
}
