package monthlyregistry

import (
	"ksef/internal/certsdb"

	"github.com/spf13/viper"
)

type InvoiceType uint8

const (
	InvoiceTypeIssued   InvoiceType = 0
	InvoiceTypeReceived InvoiceType = 1
)

type InvoiceFormCode struct {
	SystemCode    string `xml:"kodSystemowy,attr" json:"systemCode"`
	SchemaVersion string `xml:"wersjaSchemy,attr" json:"schemaVersion"`
	Value         string `xml:",chardata" json:"value"`
}

type InvoiceQRCodes struct {
	Invoice string `yaml:"invoice"`
	Offline string `yaml:"offline,omitempty"`
}

type Invoice struct {
	RefNo        string         `yaml:"ref-no"`
	KSeFRefNo    string         `yaml:"ksef-ref-no,omitempty"`
	Checksum     string         `yaml:"checksum"`
	Offline      bool           `yaml:"offline,omitempty"`
	QRCodes      InvoiceQRCodes `yaml:"qr-codes,omitempty"`
	Type         InvoiceType    `yaml:"type,omitzero"`
	UploadErrors []string       `yaml:"upload-errors,omitempty"`
}

type InvoiceMetadata struct {
	FormCode InvoiceFormCode
	Invoice  *Invoice
	Filename string
	Registry *Registry
}

type Registry struct {
	invoices []*Invoice

	dir     string
	certsDB *certsdb.CertificatesDB
	vip     *viper.Viper
}
