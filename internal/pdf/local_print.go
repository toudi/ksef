package pdf

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/pdf/cirfmf"
	"ksef/internal/pdf/printer"
	"ksef/internal/pdf/typst"
)

var ErrEngineNotConfigured = errors.New("PDF rendering engine not found in config")

func GetEngine(config *config.PDFEngineConfig) (printer.PDFPrinter, error) {
	if config.CIRFMFConfig != nil {
		return cirfmf.Printer(config.CIRFMFConfig), nil
	}

	if config.TypstConfig != nil {
		return typst.Printer(config.TypstConfig), nil
	}

	// cfg, err := config.PDFPrinterConfig(viper.GetViper())
	// if err != nil {
	// 	return nil, err
	// }

	// if cfg.LatexConfig != nil {
	// 	return latex.NewLatexPrinter(cfg.LatexConfig), nil
	// }

	// if cfg.TypstConfig != nil {
	// 	return typst.NewTypstPrinter(cfg.TypstConfig), nil
	// }

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
