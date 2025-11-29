package uploader

import "time"

type Invoice struct {
	Correction     bool      `yaml:"correction"`
	RefNo          string    `yaml:"ref-no"`
	KSeFRefNo      string    `yaml:"ksef-ref-no,omitempty"`
	Contents       string    `yaml:"contents"`
	Checksum       string    `yaml:"checksum"`
	GenerationTime time.Time `yaml:"generation-time,omitempty,omitzero"`
	Corrections    []Invoice `yaml:"corrections,omitempty"`
}
