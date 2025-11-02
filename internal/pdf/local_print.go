package pdf

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/pdf/latex"
	"ksef/internal/pdf/printer"

	"github.com/spf13/viper"
)

var ErrEngineNotConfigured = errors.New("PDF rendering engine not found in config")

func GetLocalPrintingEngine() (printer.PDFPrinter, error) {
	cfg, err := config.PDFPrinterConfig(viper.GetViper())
	if err != nil {
		return nil, err
	}

	if cfg.LatexConfig != nil {
		return latex.NewLatexPrinter(cfg.LatexConfig), nil
	}

	// engine := cfg.PDFRenderer["engine"]

	// if engine == puppeteer {
	// 	return &PuppeteerReferencePrinter{
	// 		nodeBin:         cfg.PDFRenderer["node_bin"],
	// 		browserBin:      cfg.PDFRenderer["browser_bin"],
	// 		templatePath:    cfg.PDFRenderer["template_path"],
	// 		renderingScript: cfg.PDFRenderer["rendering_script"],
	// 	}, nil
	// } else if engine == gotenberg {
	// 	return &GotenbergPrinter{
	// 		host:         cfg.PDFRenderer["host"],
	// 		templatePath: cfg.PDFRenderer["template_path"],
	// 	}, nil
	// }

	return nil, ErrEngineNotConfigured
}
