package typst

import (
	"ksef/internal/config/pdf/abstract"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

type TypstPrinterConfig struct {
	Debug     bool
	Workdir   string
	Templates abstract.Templates
}

func PrinterConfigFlags(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.Bool(cfgKeyTypstDebug, false, "tryb diagnostyczny (kopiuje wynikowy plik .typ oraz ewentualny plik z błędami do katalogu wyjściowego)")
	flags.String(cfgKeyTypstWorkdir, "/tmp", "katalog roboczy do tymczasowych plików")
	flags.String(cfgKeyTypstInvoiceTemplates, "", "ścieżka do katalogu z szablonami faktur")
	flags.String(cfgKeyTypstUpoTemplates, "", "ścieżka do katalogu z szablonem UPO")
	flags.String(cfgKeyTypstInvoiceHeaderL, "", "nagłówek faktury (strona lewa)")
	flags.String(cfgKeyTypstInvoiceHeaderC, "", "nagłówek faktury (środek)")
	flags.String(cfgKeyTypstInvoiceHeaderR, "", "nagłówek faktury (strona prawa)")
	flags.String(cfgKeyTypstInvoiceFooterL, "", "stopka faktury (strona lewa)")
	flags.String(cfgKeyTypstInvoiceFooterC, "", "stopka faktury (środek)")
	flags.String(cfgKeyTypstInvoiceFooterR, "", "stopka faktury (strona prawa)")
}

func PrinterConfig(vip *viper.Viper) (*TypstPrinterConfig, error) {
	var err error

	var config = &TypstPrinterConfig{
		Debug:   viper.GetBool(cfgKeyTypstDebug),
		Workdir: viper.GetString(cfgKeyTypstWorkdir),
	}

	config.Templates, err = parseTemplates(vip)

	return config, err
}
