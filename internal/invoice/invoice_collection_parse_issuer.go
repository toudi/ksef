package invoice

import (
	"encoding/xml"
	"fmt"
	"os"
)

type XMLInvoice struct {
	XMLName       xml.Name `xml:"Faktura"`
	Issuer        string   `xml:"Podmiot1>DaneIdentyfikacyjne>NIP"`
	InvoiceNumber string   `xml:"Fa>P_2"`
}

func parseInvoiceIssuer(sourceFile string) (string, error) {
	invoice, err := ParseInvoice(sourceFile)
	if err != nil {
		return "", fmt.Errorf("cannot parse issuer: %v", err)
	}
	return invoice.Issuer, nil
}

func ParseInvoice(sourceFile string) (*XMLInvoice, error) {
	var invoice XMLInvoice
	xmlContents, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read invoice file: %v", err)
	}
	if err = xml.Unmarshal(xmlContents, &invoice); err != nil {
		return nil, fmt.Errorf("unable to parse xml invoice: %v", err)
	}

	return &invoice, nil
}
