package latex

import (
	"ksef/internal/config/pdf/abstract"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyLatexDebug            = "pdf.latex.debug"
	cfgKeyLatexWorkdir          = "pdf.latex.workdir"
	cfgKeyLatexDocker           = "pdf.latex.docker"
	cfgKeyLatexPodman           = "pdf.latex.podman"
	cfgKeyLatexInvoiceTemplates = "pdf.latex.invoice.templates"
	cfgKeyLatexInvoiceHeaderL   = "pdf.latex.invoice.header.left"
	cfgKeyLatexInvoiceHeaderC   = "pdf.latex.invoice.header.center"
	cfgKeyLatexInvoiceHeaderR   = "pdf.latex.invoice.header.right"
	cfgKeyLatexInvoiceFooterL   = "pdf.latex.invoice.footer.left"
	cfgKeyLatexInvoiceFooterC   = "pdf.latex.invoice.footer.center"
	cfgKeyLatexInvoiceFooterR   = "pdf.latex.invoice.footer.right"
	cfgKeyLatexUpoTemplates     = "pdf.latex.upo.templates"
)

type LatexPrinterConfig struct {
	Debug     bool
	Workdir   string
	Docker    bool
	Podman    bool
	Templates abstract.Templates
}

func PrinterConfigFlags(cmd *cobra.Command, flags *pflag.FlagSet) {
	flags.Bool(cfgKeyLatexDebug, false, "tryb diagnostyczny (kopiuje wynikowy plik .tex oraz ewentualny plik z błędami do katalogu wyjściowego)")
	flags.String(cfgKeyLatexWorkdir, "/tmp", "katalog roboczy do tymczasowych plików")
	flags.Bool(cfgKeyLatexDocker, true, "użyj LaTeX dostarczonego przez kontener docker")
	flags.Bool(cfgKeyLatexPodman, false, "użyj LaTeX dostarczonego przez kontener podman")
	flags.String(cfgKeyLatexInvoiceTemplates, "", "ścieżka do katalogu z szablonami faktur")
	flags.String(cfgKeyLatexUpoTemplates, "", "ścieżka do katalogu z szablonem UPO")
	flags.String(cfgKeyLatexInvoiceHeaderL, "", "nagłówek faktury (strona lewa)")
	flags.String(cfgKeyLatexInvoiceHeaderC, "", "nagłówek faktury (środek)")
	flags.String(cfgKeyLatexInvoiceHeaderR, "", "nagłówek faktury (strona prawa)")
	flags.String(cfgKeyLatexInvoiceFooterL, "", "stopka faktury (strona lewa)")
	flags.String(cfgKeyLatexInvoiceFooterC, "", "stopka faktury (środek)")
	flags.String(cfgKeyLatexInvoiceFooterR, "", "stopka faktury (strona prawa)")

	cmd.MarkFlagsMutuallyExclusive(cfgKeyLatexDocker, cfgKeyLatexPodman)
}

func PrinterConfig(vip *viper.Viper) (*LatexPrinterConfig, error) {
	var err error

	var config = &LatexPrinterConfig{
		Debug:   viper.GetBool(cfgKeyLatexDebug),
		Workdir: viper.GetString(cfgKeyLatexWorkdir),
		Docker:  viper.GetBool(cfgKeyLatexDocker),
		Podman:  viper.GetBool(cfgKeyLatexPodman),
	}

	config.Templates, err = parseTemplates(vip)

	return config, err
}
