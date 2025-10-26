package config

import (
	"errors"
	"ksef/internal/config/pdf/latex"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	printerEngineLatex     = "latex"
	printerEnginePuppeteer = "puppeteer"
	printerEngineGotenberg = "gotenberg"

	cfgKeyPdfEngine string = "pdf.engine"
)

var (
	errInvalidEngine = errors.New("nieprawidłowa wartość opcji pdf.engine")
)

type pdfPrinterConfig struct {
	LatexConfig *latex.LatexPrinterConfig
}

func PDFPrinterConfig(vip *viper.Viper) (config pdfPrinterConfig, err error) {
	var engine = vip.GetString(cfgKeyPdfEngine)

	switch engine {
	case printerEngineLatex:
		var latexConfig *latex.LatexPrinterConfig
		latexConfig, err = latex.PrinterConfig(vip)
		if err != nil {
			break
		}
		config.LatexConfig = latexConfig
	default:
		err = errInvalidEngine
	}

	return config, err
}

func PDFPrinterFlags(cmd *cobra.Command, flags *pflag.FlagSet) error {
	flags.String(cfgKeyPdfEngine, "", "silnik renderujący")
	latex.PrinterConfigFlags(cmd, flags)
	return nil
}
