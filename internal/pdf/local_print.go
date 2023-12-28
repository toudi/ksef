package pdf

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/registry"
)

const (
	puppeteer = "puppeteer"
	gotenberg = "gotenberg"
)

type PDFPrinter interface {
	Print(contentBase64 string, meta registry.Invoice, output string) error
}

var ErrEngineNotConfigured = errors.New("PDF rendering engine not found in config")

func GetLocalPrintingEngine() (PDFPrinter, error) {
	engine := config.Config.PDFRenderer["engine"]

	if engine == puppeteer {
		return &PuppeteerReferencePrinter{
			nodeBin:         config.Config.PDFRenderer["node_bin"],
			browserBin:      config.Config.PDFRenderer["browser_bin"],
			templatePath:    config.Config.PDFRenderer["template_path"],
			renderingScript: config.Config.PDFRenderer["rendering_script"],
		}, nil
	} else if engine == gotenberg {
		return &GotenbergPrinter{
			host:         config.Config.PDFRenderer["host"],
			templatePath: config.Config.PDFRenderer["template_path"],
		}, nil
	}

	return nil, ErrEngineNotConfigured
}
