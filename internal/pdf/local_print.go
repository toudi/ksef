package pdf

import (
	"errors"
	"fmt"
	"ksef/internal/config"
	"ksef/internal/logging"
	"ksef/internal/registry"
	"os"
	"os/exec"
	"path"
	"strings"
)

const puppeteer = "puppeteer"

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
	}

	return nil, ErrEngineNotConfigured
}

type PuppeteerReferencePrinter struct {
	nodeBin         string
	browserBin      string
	templatePath    string
	renderingScript string
}

func (p *PuppeteerReferencePrinter) Print(
	contentBase64 string,
	invoiceMeta registry.Invoice,
	output string,
) error {
	replacer := strings.NewReplacer(
		"__qrcode_url__", invoiceMeta.SEIQRCode,
		"__invoice_sei_ref_no__", invoiceMeta.SEIReferenceNumber,
		"__invoice_base64__", contentBase64,
	)

	templateDir := path.Dir(p.templatePath)
	templateContent, err := os.ReadFile(p.templatePath)
	if err != nil {
		return fmt.Errorf("unable to read template file: %v", err)
	}

	renderFilePath := path.Join(templateDir, "render.html")
	destFile, err := os.Create(renderFilePath)
	if err != nil {
		return fmt.Errorf("unable to create pre-rendered file: %v", err)
	}

	_, err = replacer.WriteString(destFile, string(templateContent))
	destFile.Close()
	if err != nil {
		return fmt.Errorf("unable to write to destFile: %v", err)
	}

	command := p.nodeBin
	commandArgs := []string{
		p.renderingScript,
		p.browserBin,
		renderFilePath,
		output,
	}

	logging.PDFRendererLogger.Debug("local PDF render", "command", command, "args", commandArgs)

	cmd := exec.Command(command, commandArgs...)
	b := new(strings.Builder)
	cmd.Stderr = b

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error from stderr: %v\n", b)
		return fmt.Errorf("error during conversion to PDF: %v", err)
	}

	return nil
}
