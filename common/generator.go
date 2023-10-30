package common

import (
	"ksef/common/xml"
)

type Generator interface {
	InvoiceToXMLTree(invoice *Invoice) (*xml.Node, error)
	LineHandler(invoice *Invoice, section string, data map[string]string, invoiceReady func() error) error
	// method that returns issuer tax identification number from the invoice object
	IssuerTIN() string
}

var SectionInvoice = "faktura.fa"
