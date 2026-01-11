package pdf

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/pdf/cirfmf"
	pdfconfig "ksef/internal/pdf/config"
	"ksef/internal/pdf/printer"
	"ksef/internal/pdf/typst"

	"github.com/spf13/viper"
)

var ErrEngineNotConfigured = errors.New("PDF rendering engine not found in config")

func GetEngine(config *pdfconfig.PDFEngineConfig) (printer.PDFPrinter, error) {
	if config.CIRFMFConfig != nil {
		return cirfmf.Printer(config.CIRFMFConfig), nil
	}

	if config.TypstConfig != nil {
		return typst.Printer(config.TypstConfig), nil
	}
	return nil, ErrEngineNotConfigured
}

func GetUPOPrinter(vip *viper.Viper) (printer.PDFPrinter, error) {
	pdfConfig, err := config.GetPDFPrinterConfig(vip)
	if err != nil {
		return nil, err
	}

	engineConfig, err := pdfConfig.GetEngine(pdfconfig.UsageSelector{Usage: "upo"})
	if err != nil {
		return nil, err
	}

	return GetEngine(engineConfig)
}

func GetInvoicePrinter(vip *viper.Viper, usage pdfconfig.UsageSelector) (printer.PDFPrinter, error) {
	pdfConfig, err := config.GetPDFPrinterConfig(vip)
	if err != nil {
		return nil, err
	}

	engineConfig, err := pdfConfig.GetEngine(usage)
	if err != nil {
		return nil, err
	}

	return GetEngine(engineConfig)
}
