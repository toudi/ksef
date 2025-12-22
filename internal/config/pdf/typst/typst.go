package typst

import (
	"github.com/spf13/viper"
)

const (
	cfgKeyTypstDebug            = "pdf.typst.debug"
	cfgKeyTypstWorkdir          = "pdf.typst.workdir"
	cfgKeyTypstInvoiceTemplates = "pdf.typst.invoice.templates"
	cfgKeyTypstInvoiceHeaderL   = "pdf.typst.invoice.header.left"
	cfgKeyTypstInvoiceHeaderC   = "pdf.typst.invoice.header.center"
	cfgKeyTypstInvoiceHeaderR   = "pdf.typst.invoice.header.right"
	cfgKeyTypstInvoiceFooterL   = "pdf.typst.invoice.footer.left"
	cfgKeyTypstInvoiceFooterC   = "pdf.typst.invoice.footer.center"
	cfgKeyTypstInvoiceFooterR   = "pdf.typst.invoice.footer.right"
	cfgKeyTypstUpoTemplates     = "pdf.typst.upo.templates"
)

type HeaderFooterConfig struct {
	Left   string `yaml:"left"`
	Center string `yaml:"center"`
	Right  string `yaml:"right"`
}

type TypstInvoicePrinterConfig struct {
	Template string             `yaml:"template"`
	Header   HeaderFooterConfig `yaml:"header"`
	Footer   HeaderFooterConfig `yaml:"footer"`
}

type TypstUPOPrinterConfig struct {
	Template string `yaml:"template"`
}

type TypstPrinterConfig struct {
	Debug     bool                      `yaml:"debug"`
	Workdir   string                    `yaml:"workdir"`
	Templates string                    `yaml:"templates-dir"`
	Invoice   TypstInvoicePrinterConfig `yaml:"invoice"`
	UPO       TypstUPOPrinterConfig     `yaml:"upo"`
}

func PrinterConfig(vip *viper.Viper) (*TypstPrinterConfig, error) {
	var err error

	config := &TypstPrinterConfig{
		Debug:   viper.GetBool(cfgKeyTypstDebug),
		Workdir: viper.GetString(cfgKeyTypstWorkdir),
	}

	// config.Templates, err = parseTemplates(vip)

	return config, err
}
