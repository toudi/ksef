package monthlyregistry

import (
	"ksef/internal/certsdb"

	sessionTypes "ksef/internal/client/v2/session/types"

	"github.com/spf13/viper"
)

type InvoiceType uint8

const (
	InvoiceTypeIssued   InvoiceType = 0
	InvoiceTypeReceived InvoiceType = 1
)

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
	FormCode sessionTypes.InvoiceFormCode
	Invoice  *Invoice
	Filename string
}

type UploadSession struct {
	RefNo     string     `yaml:"ref-no"`
	Processed bool       `yaml:"processed,omitempty"`
	Invoices  []*Invoice `yaml:"invoices"`
}

type Registry struct {
	invoices       []*Invoice       `yaml:"invoices"`
	uploadSessions []*UploadSession `yaml:"upload-sessions,omitempty"`

	dir     string
	certsDB *certsdb.CertificatesDB
	vip     *viper.Viper
}
