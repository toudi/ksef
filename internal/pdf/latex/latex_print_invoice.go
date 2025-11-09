package latex

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"ksef/internal/config/pdf/latex"
	"ksef/internal/pdf/printer"
	"ksef/internal/registry/types"
	"ksef/internal/sei/generators/fa_3_1"
)

func NewLatexPrinter(cfg *latex.LatexPrinterConfig) *LatexPrinter {
	return &LatexPrinter{cfg: cfg}
}

type LatexPrinter struct {
	cfg *latex.LatexPrinterConfig
}

type LatexTemplateVars struct {
	Invoice  *printer.Invoice
	Registry types.Invoice
	Qrcodes  types.InvoiceQRCodes
	Header   latex.HeaderFooterSettings
	Footer   latex.HeaderFooterSettings
}

func (lp *LatexPrinter) Print(contentBase64 string, meta types.Invoice, output string) error {
	var invoiceXMLBuffer bytes.Buffer
	invoiceXMLBytes, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return err
	}
	if _, err = invoiceXMLBuffer.Write(invoiceXMLBytes); err != nil {
		return err
	}

	var i = printer.Invoice{ArrayElements: fa_3_1.FA_3_1ArrayElements}
	if err = xml.Unmarshal(invoiceXMLBytes, &i); err != nil {
		return err
	}

	var templateData = LatexTemplateVars{
		Invoice:  &i,
		Registry: meta,
		Qrcodes:  meta.QRCodes,
		Header:   lp.cfg.Templates.Invoice.Header,
		Footer:   lp.cfg.Templates.Invoice.Footer,
	}

	return lp.printTemplate(lp.cfg.Templates.Invoice.Template, "invoice", templateData, output)
}
