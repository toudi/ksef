package monthlyregistry

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"time"

	sessionTypes "ksef/internal/client/v2/session/types"

	"github.com/spf13/viper"
)

type InvoiceType uint8

const (
	InvoiceTypeIssued     InvoiceType = 0
	InvoiceTypeReceived   InvoiceType = 1
	InvoiceTypePayer      InvoiceType = 2
	InvoiceTypeAuthorized InvoiceType = 3
)

var invoiceTypeToPrinterUsage = map[InvoiceType]string{
	InvoiceTypeIssued:     "issued",
	InvoiceTypeReceived:   "received",
	InvoiceTypePayer:      "payer",
	InvoiceTypeAuthorized: "authorized",
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
	PrintoutData map[string]any `yaml:"printout-data,omitempty"`
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

type SyncParams struct {
	LastTimestamp time.Time              `yaml:"last-timestamp,omitempty"`
	SubjectTypes  []invoices.SubjectType `yaml:"subject-types,omitempty"`
}

type Registry struct {
	Invoices   []*Invoice  `yaml:"invoices"`
	SyncParams *SyncParams `yaml:"sync,omitempty"`

	dir     string
	certsDB *certsdb.CertificatesDB
	vip     *viper.Viper
}
