package uploader

import "time"

type Invoice struct {
	RefNo          string    `yaml:"ref-no"`
	KSeFRefNo      string    `yaml:"ksef-ref-no"`
	Contents       string    `yaml:"contents"`
	Checksum       string    `yaml:"checksum"`
	GenerationTime time.Time `yaml:"generation-time"`
	Corrections    []Invoice `yaml:"corrections,omitempty"`
}
