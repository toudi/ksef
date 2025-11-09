package types

type InvoiceQRCodes struct {
	Invoice     string `yaml:"invoice,omitempty"`
	Certificate string `yaml:"certificate,omitempty"`
}

type Invoice struct {
	ReferenceNumber     string         `yaml:"referenceNumber,omitempty"`
	KSeFReferenceNumber string         `yaml:"ksefReferenceNumber,omitempty"`
	QRCodes             InvoiceQRCodes `yaml:"qrcodes"`
	IssueDate           string         `yaml:"issueDate,omitempty"`
	SubjectFrom         InvoiceSubject `yaml:"subjectFrom,omitempty"`
	SubjectTo           InvoiceSubject `yaml:"subjectTo,omitempty"`
	InvoiceType         string         `yaml:"invoiceType,omitempty"`
	Net                 string         `yaml:"net,omitempty"`
	Vat                 string         `yaml:"vat,omitempty"`
	Gross               string         `yaml:"gross,omitempty"`
	Checksum            string         `yaml:"checksum,omitempty"`
	Offline             bool           `yaml:"offline,omitempty"`
}
