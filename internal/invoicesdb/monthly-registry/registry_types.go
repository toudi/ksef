package monthlyregistry

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/invoicesdb/shared"
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

type InvoiceIssuer struct {
	NIP  string `yaml:"nip"`
	Name string `yaml:"name"`
}

// Annotations is a slice of Annotation rules that marshals as a flat YAML list
// under the "annotations" key (no intermediate "item-rules" field).
type Annotations []shared.Annotation

func (a Annotations) MarshalYAML() (any, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return []shared.Annotation(a), nil
}

func (a *Annotations) UnmarshalYAML(unmarshal func(any) error) error {
	var rules []shared.Annotation
	if err := unmarshal(&rules); err != nil {
		return err
	}
	*a = rules
	return nil
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
	Issuer       *InvoiceIssuer `yaml:"issuer,omitempty"`
	OrdNum       int            `yaml:"ord-num"`
	Annotations  Annotations    `yaml:"annotations,omitempty"`
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

type OrdNumState struct {
	InvoiceType InvoiceType `yaml:"invoice-type"`
	Count       int         `yaml:"count"`
}

type OrdNums []OrdNumState

func (o OrdNums) ToMap() map[InvoiceType]int {
	result := make(map[InvoiceType]int)
	for _, entry := range o {
		result[entry.InvoiceType] = entry.Count
	}
	return result
}

type OrdNumsMap map[InvoiceType]int

func (o OrdNumsMap) ToSlice() (result OrdNums) {
	for invoiceType, count := range o {
		result = append(result, OrdNumState{InvoiceType: invoiceType, Count: count})
	}
	return result
}

type Registry struct {
	Invoices     []*Invoice  `yaml:"invoices"`
	SyncParams   *SyncParams `yaml:"sync,omitempty"`
	OrdNums      OrdNumsMap  `yaml:"-"`
	SavedOrdNums OrdNums     `yaml:"ord-nums,omitempty"`

	checksumIndex map[string]int          `yaml:"-"`
	dir           string                  `yaml:"-"`
	certsDB       *certsdb.CertificatesDB `yaml:"-"`
	vip           *viper.Viper            `yaml:"-"`
}
