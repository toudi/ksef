package monthly

import (
	"ksef/core/ksef"
	"ksef/internal/invoicesdb/shared"
	"time"
)

type InvoiceType uint8

const (
	InvoiceTypeIssued     InvoiceType = 0
	InvoiceTypeReceived   InvoiceType = 1
	InvoiceTypePayer      InvoiceType = 2
	InvoiceTypeAuthorized InvoiceType = 3
)

var subjectTypeToInvoiceType = map[ksef.SubjectType]InvoiceType{
	ksef.SubjectTypeIssuer:     InvoiceTypeIssued,
	ksef.SubjectTypeRecipient:  InvoiceTypeReceived,
	ksef.SubjectTypePayer:      InvoiceTypePayer,
	ksef.SubjectTypeAuthorized: InvoiceTypeAuthorized,
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

// SyncParams mirrors the legacy registry's sync section exactly.
type SyncParams struct {
	LastTimestamp time.Time          `yaml:"last-timestamp,omitempty"`
	SubjectTypes  []ksef.SubjectType `yaml:"subject-types,omitempty"`
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

func ordNumsToSlice(input map[InvoiceType]int) OrdNums {
	result := make(OrdNums, 0, len(input))
	for invoiceType, count := range input {
		result = append(result, OrdNumState{InvoiceType: invoiceType, Count: count})
	}
	return result
}

type state struct {
	Invoices     []*Invoice  `yaml:"invoices"`
	SyncParams   *SyncParams `yaml:"sync,omitempty"`
	SavedOrdNums OrdNums     `yaml:"ord-nums,omitempty"`
}
