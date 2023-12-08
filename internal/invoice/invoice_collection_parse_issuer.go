package invoice

import (
	"encoding/xml"
	"fmt"
	"os"
)

type xmlInvoice struct {
	XMLName xml.Name `xml:"Faktura"`
	Issuer  string   `xml:"Podmiot1>DaneIdentyfikacyjne>NIP"`
}

func parseInvoiceIssuer(sourceFile string) (string, error) {
	var invoice xmlInvoice
	xmlContents, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", fmt.Errorf("unable to read invoice file: %v", err)
	}
	if err = xml.Unmarshal(xmlContents, &invoice); err != nil {
		return "", fmt.Errorf("unable to parse xml invoice: %v", err)
	}
	return invoice.Issuer, nil
}
