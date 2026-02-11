package annualregistry

import "time"

type Invoice struct {
	Correction     bool       `yaml:"correction,omitempty"`
	RefNo          string     `yaml:"ref-no"`
	KSeFRefNo      string     `yaml:"ksef-ref-no,omitempty"`
	Contents       string     `yaml:"contents,omitempty"` // base64 of gzip of cbor of the raw invoice data
	Checksum       string     `yaml:"checksum"`           // sha256 checksum of the generated XML document - this is **NOT** the same as ContentChecksum
	GenerationTime time.Time  `yaml:"generation-time,omitempty,omitzero"`
	Corrections    []*Invoice `yaml:"corrections,omitempty"`
}
