package metadata

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Invoice struct {
	XMLName xml.Name `xml:"Faktura"`
	Issuer  string   `xml:"Podmiot1>DaneIdentyfikacyjne>NIP"`
}

func ParseIssuerFromInvoice(sourceFile string) (string, error) {
	var invoice Invoice
	xmlContents, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", fmt.Errorf("unable to read invoice file: %v", err)
	}
	if err = xml.Unmarshal(xmlContents, &invoice); err != nil {
		return "", fmt.Errorf("unable to parse xml invoice: %v", err)
	}
	return invoice.Issuer, nil
}
