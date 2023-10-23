package common

import (
	"ksef/common/xml"
	"ksef/metadata"
)

type Generator interface {
	InvoiceToXMLTree(invoice *Invoice) (*xml.Node, error)
	LineHandler(string, map[string]string) error
	Save(string) error
	// method that returns issuer tax identification number from the invoice object
	IssuerTIN() string
	PopulateMetadata(meta *metadata.Metadata, sourceFile string) error
}

var SectionInvoice = "faktura.fa"
