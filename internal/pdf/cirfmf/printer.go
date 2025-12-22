package cirfmf

import (
	"ksef/internal/config/pdf/cirfmf"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/pdf/printer"
	"os/exec"
	"path/filepath"
)

type cirfmfPrinter struct {
	config cirfmf.PrinterConfig
}

func Printer(config *cirfmf.PrinterConfig) printer.PDFPrinter {
	return &cirfmfPrinter{
		config: *config,
	}
}

func (p *cirfmfPrinter) PrintUPO(srcFile string, output string) error {
	cmd := exec.Command(
		p.config.NodeBinPath,
		filepath.Join(p.config.TemplatesDir, "print.js"),
		"upo",
		srcFile,
		output,
	)
	return cmd.Run()
}

func (p *cirfmfPrinter) PrintInvoice(
	srcFile string,
	output string,
	meta *monthlyregistry.InvoicePrintingMeta,
) error {
	cmd := exec.Command(
		p.config.NodeBinPath,
		filepath.Join(p.config.TemplatesDir, "print.js"),
		"invoice",
		srcFile,
		output,
		meta.Invoice.KSeFRefNo,
		meta.Invoice.QRCodes.Invoice,
	)
	return cmd.Run()
}
