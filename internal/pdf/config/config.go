package config

import (
	"ksef/internal/config/pdf/cirfmf"
	"ksef/internal/config/pdf/typst"
)

type UsageSelector struct {
	Usage        string
	Participants map[string]any // only applicable to invoices
}

type PDFEngineConfig struct {
	UsageRaw     any                       `yaml:"usage"`
	Condition    string                    `yaml:"if,omitempty"`
	TypstConfig  *typst.TypstPrinterConfig `yaml:"typst,omitempty"`
	CIRFMFConfig *cirfmf.PrinterConfig     `yaml:"cirfmf,omitempty"`
}

func (c PDFEngineConfig) Usage() []string {
	if usageSlice, ok := c.UsageRaw.([]any); ok {
		var usageStringSlice []string
		for _, usageStr := range usageSlice {
			usageStringSlice = append(usageStringSlice, usageStr.(string))
		}
		return usageStringSlice
	} else if usageStr, ok := c.UsageRaw.(string); ok {
		return []string{usageStr}
	}
	return []string{}
}
