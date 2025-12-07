package monthlyregistry

import (
	"encoding/xml"
	"errors"
	"os"
)

var (
	errUnableToReadFile = errors.New("unable to read invoice file")
	errUnableToParseXML = errors.New("unable to parse XML")
)

type XMLInvoice struct {
	XMLName        xml.Name        `xml:"Faktura"`
	HeaderFormCode InvoiceFormCode `xml:"Naglowek>KodFormularza"`
	Issuer         string          `xml:"Podmiot1>DaneIdentyfikacyjne>NIP"`
	Issued         string          `xml:"Fa>P_1"`
	InvoiceNumber  string          `xml:"Fa>P_2"`
}

func (r *Registry) getInvoiceMetadata(input *Invoice, ordNo int) (*InvoiceMetadata, error) {
	invoiceFilename := r.getIssuedInvoiceFilename(input.RefNo, ordNo)
	// in order to obtain the metadata we need to parse XML and extract the formcode
	if xmlInvoice, err := parseInvoice(invoiceFilename); err != nil {
		return nil, err
	} else {
		return &InvoiceMetadata{
			FormCode: xmlInvoice.HeaderFormCode,
			Invoice:  input,
			Filename: invoiceFilename,
			Registry: r, // we add pointer to the registry so that the uploader can later find the invoice and save it's reference number
		}, nil
	}
}

func parseInvoice(sourceFile string) (*XMLInvoice, error) {
	var invoice XMLInvoice
	xmlContents, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, errors.Join(errUnableToReadFile, err)
	}
	if err = xml.Unmarshal(xmlContents, &invoice); err != nil {
		return nil, errors.Join(errUnableToParseXML, err)
	}

	return &invoice, nil
}
