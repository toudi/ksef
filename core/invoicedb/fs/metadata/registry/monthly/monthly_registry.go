package monthly

import "ksef/internal/invoicesdb/shared"

type InvoiceType uint8

const (
	InvoiceTypeIssued     InvoiceType = 0
	InvoiceTypeReceived   InvoiceType = 1
	InvoiceTypePayer      InvoiceType = 2
	InvoiceTypeAuthorized InvoiceType = 3
)

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

type state struct {
	Invoices []*Invoice `yaml:"invoices"`
}
