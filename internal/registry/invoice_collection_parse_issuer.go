package registry

import (
	"encoding/xml"
	"fmt"
	"os"
)

type XMLInvoice struct {
	XMLName        xml.Name        `xml:"Faktura"`
	HeaderFormCode InvoiceFormCode `xml:"Naglowek>KodFormularza"`
	Issuer         string          `xml:"Podmiot1>DaneIdentyfikacyjne>NIP"`
	Issued         string          `xml:"Fa>P_1"`
	InvoiceNumber  string          `xml:"Fa>P_2"`
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

func ReadAndParseInvoice(sourceFile string) (*XMLInvoice, []byte, error) {
	var invoice XMLInvoice
	xmlContents, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read invoice file: %v", err)
	}
	if err = xml.Unmarshal(xmlContents, &invoice); err != nil {
		return nil, nil, fmt.Errorf("unable to parse xml invoice: %v", err)
	}

	return &invoice, xmlContents, nil
}
