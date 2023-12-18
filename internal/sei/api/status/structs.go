package status

import "errors"

type KsefInvoiceIdType struct {
	InvoiceNumber          string `xml:"NumerFaktury" json:"invoiceNumber" yaml:"invoiceNumber"`
	KSeFInvoiceReferenceNo string `xml:"NumerKSeFDokumentu" json:"ksefDocumentId" yaml:"ksefDocumentId"`
	DocumentChecksum       string `xml:"SkrotDokumentu"`
}

type StatusInfo struct {
	SelectedFormat string              `json:"-" yaml:"-"`
	SourcePath     string              `json:"-" yaml:"-"`
	Environment    string              `json:"env" yaml:"env"`
	SessionID      string              `json:"sessionId" yaml:"sessionId"`
	Issuer         string              `json:"issuer" yaml:"issuer"`
	InvoiceIds     []KsefInvoiceIdType `json:"invoiceIds,omitempty" yaml:"invoiceIds,omitempty"`
}

func (s *StatusInfo) GetSEIRefNo(invoiceNo string) (string, error) {
	for _, invoice := range s.InvoiceIds {
		if invoice.InvoiceNumber == invoiceNo {
			return invoice.KSeFInvoiceReferenceNo, nil
		}
	}

	return "", errors.New("invoice number could not be found")
}
