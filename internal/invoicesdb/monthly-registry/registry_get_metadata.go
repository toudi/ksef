package monthlyregistry

import (
	"encoding/xml"
	"errors"
	sessionTypes "ksef/internal/client/v2/session/types"
	"ksef/internal/utils"
	"os"
	"time"
)

var (
	errUnableToReadFile = errors.New("unable to read invoice file")
	errUnableToParseXML = errors.New("unable to parse XML")
)

type Item struct {
	OrdNo   string `xml:"NrWierszaFa"`
	Name    string `xml:"P_7"`
	Index   string `xml:"Indeks"`
	GTIN    string `xml:"GTIN"`
	PKWiU   string `xml:"PKWiU"`
	CN      string `xml:"CN"`
	PKOB    string `xml:"PKOB"`
	VATRate string `xml:"P_12"`
}

type XMLInvoice struct {
	XMLName             xml.Name                     `xml:"Faktura"`
	HeaderFormCode      sessionTypes.InvoiceFormCode `xml:"Naglowek>KodFormularza"`
	GeneratedTime       time.Time                    `xml:"Naglowek>DataWytworzeniaFa"`
	Issuer              string                       `xml:"Podmiot1>DaneIdentyfikacyjne>NIP"`
	IssuerName          string                       `xml:"Podmiot1>DaneIdentyfikacyjne>Nazwa"`
	RecipientName       string                       `xml:"Podmiot2>DaneIdentyfikacyjne>Nazwa"`
	Recipient           string                       `xml:"Podmiot2>DaneIdentyfikacyjne>NIP"`
	RecipientVatIdEU    string                       `xml:"Podmiot2>DaneIdentyfikacyjne>NrVatUE"`
	RecipientVatIdNonEU string                       `xml:"Podmiot2>DaneIdentyfikacyjne>NrID"`
	Issued              string                       `xml:"Fa>P_1"`
	InvoiceNumber       string                       `xml:"Fa>P_2"`
	Items               []Item                       `xml:"Fa>FaWiersz"`
}

func (xi *XMLInvoice) Participants() map[string]any {
	sellerMap := map[string]string{
		"nip": xi.Issuer,
	}
	buyerMap := map[string]any{}
	if xi.Recipient != "" {
		buyerMap["nip"] = xi.Recipient
	} else if xi.RecipientVatIdEU != "" {
		buyerMap["nr_vat_ue"] = xi.RecipientVatIdEU
	} else if xi.RecipientVatIdNonEU != "" {
		buyerMap["nr_id"] = xi.RecipientVatIdNonEU
	}
	return map[string]any{
		"seller": sellerMap,
		"buyer":  buyerMap,
	}
}

func (r *Registry) getInvoiceMetadata(input *Invoice, ordNo int) (*InvoiceMetadata, error) {
	invoiceFilename := r.getIssuedInvoiceFilename(input.RefNo, ordNo)
	// in order to obtain the metadata we need to parse XML and extract the formcode
	if xmlInvoice, _, err := ParseInvoice(invoiceFilename); err != nil {
		return nil, err
	} else {
		return &InvoiceMetadata{
			FormCode: xmlInvoice.HeaderFormCode,
			Invoice:  input,
			Filename: invoiceFilename,
		}, nil
	}
}

func ParseInvoice(sourceFile string) (*XMLInvoice, string, error) {
	var invoice XMLInvoice
	xmlContents, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, "", errors.Join(errUnableToReadFile, err)
	}
	if err = xml.Unmarshal(xmlContents, &invoice); err != nil {
		return nil, "", errors.Join(errUnableToParseXML, err)
	}

	return &invoice, utils.Sha256Hex(xmlContents), nil
}

func InvoiceParticipants(sourceFile string) (map[string]any, error) {
	inv, _, err := ParseInvoice(sourceFile)
	if err != nil {
		return nil, err
	}
	return inv.Participants(), nil
}
