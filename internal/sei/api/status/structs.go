package status

type KsefInvoiceIdType struct {
	InvoiceNumber          string `xml:"NumerFaktury" json:"invoiceNumber" yaml:"invoiceNumber"`
	KSeFInvoiceReferenceNo string `xml:"NumerKSeFDokumentu" json:"ksefDocumentId" yaml:"ksefDocumentId"`
}

type StatusInfo struct {
	SelectedFormat string              `json:"-" yaml:"-"`
	SourcePath     string              `json:"-" yaml:"-"`
	Environment    string              `json:"env" yaml:"env"`
	SessionID      string              `json:"sessionId" yaml:"sessionId"`
	Issuer         string              `json:"issuer" yaml:"issuer"`
	InvoiceIds     []KsefInvoiceIdType `json:"invoiceIds,omitempty" yaml:"invoiceIds,omitempty"`
}
