package config

import (
	"errors"
	"ksef/internal/config/pdf/latex"
	"ksef/internal/config/pdf/typst"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	printerEngineLatex     = "latex"
	printerEngineTypst     = "typst"
	printerEnginePuppeteer = "puppeteer"
	printerEngineGotenberg = "gotenberg"

	cfgKeyPdfEngine string = "pdf.engine"
)

var (
	errInvalidEngine = errors.New("nieprawidłowa wartość opcji pdf.engine")
)

type pdfPrinterConfig struct {
	LatexConfig *latex.LatexPrinterConfig
	TypstConfig *typst.TypstPrinterConfig
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
	case printerEngineTypst:
		var typstConfig *typst.TypstPrinterConfig
		typstConfig, err = typst.PrinterConfig(vip)
		if err != nil {
			break
		}
		config.TypstConfig = typstConfig
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
