package invoice

import "ksef/internal/invoicesdb/shared"

type InvoiceType uint8

const (
	InvoiceTypeIssued     InvoiceType = 0
	InvoiceTypeReceived   InvoiceType = 1
	InvoiceTypePayer      InvoiceType = 2
	InvoiceTypeAuthorized InvoiceType = 3
)

type InvoiceQRCodes struct {
	Invoice string
	Offline string
}

type InvoiceIssuer struct {
	NIP  string
	Name string
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
	RefNo        string
	KSeFRefNo    string
	Checksum     string
	Offline      bool
	QRCodes      InvoiceQRCodes
	Type         InvoiceType
	UploadErrors []string
	PrintoutData map[string]any
	Issuer       *InvoiceIssuer
	OrdNum       int
	Annotations  Annotations
}
