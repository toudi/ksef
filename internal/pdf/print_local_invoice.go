package pdf

import (
	"encoding/base64"
	"ksef/internal/pdf/printer"
	"ksef/internal/registry"
	"os"
	"strings"
)

func PrintLocalInvoice(engine printer.PDFPrinter, invoice registry.Invoice, filename string) error {
	var base64Encoder = base64.StdEncoding
	invoiceContents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	invoiceBase64 := base64Encoder.EncodeToString(invoiceContents)
	return engine.Print(
		invoiceBase64,
		invoice,
		strings.Replace(filename, ".xml", ".pdf", 1),
	)
}
