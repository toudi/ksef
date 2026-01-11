package monthlyregistry

import (
	"path/filepath"
)

type InvoiceMeta struct {
	KSeFRefNo string         `yaml:"ksef-ref-no"`
	QRCodes   InvoiceQRCodes `yaml:"qr-codes"`
}

type HeaderFooter struct {
	Left   string `yaml:"left"`
	Center string `yaml:"center"`
	Right  string `yaml:"right"`
}

type PageSettings struct {
	Header HeaderFooter `yaml:"header"`
	Footer HeaderFooter `yaml:"footer"`
}

type InvoicePrintingMeta struct {
	Usage        string         `yaml:"-"` // entry in the config usage slice
	Invoice      InvoiceMeta    `yaml:"invoice"`
	Page         PageSettings   `yaml:"page"`
	Printout     map[string]any `yaml:"printout"`
	Participants map[string]any `yaml:"-"`
}

func GetInvoicePrintingMeta(srcFile string) (*InvoicePrintingMeta, error) {
	// for the invoice, we need to extract meta information from the monthly registry
	registryDir := filepath.Join(filepath.Dir(srcFile), "..")
	monthlyRegistry, err := Open(registryDir)
	if err != nil {
		return nil, err
	}
	xmlInvoice, checksum, err := ParseInvoice(srcFile)
	if err != nil {
		return nil, err
	}
	invoice := monthlyRegistry.GetInvoiceByChecksum(checksum)
	if invoice == nil {
		return nil, errUnableToFindInvoice
	}

	usage := "invoice:" + invoiceTypeToPrinterUsage[invoice.Type]

	return &InvoicePrintingMeta{
		Usage: usage,
		Invoice: InvoiceMeta{
			KSeFRefNo: invoice.KSeFRefNo,
			QRCodes:   invoice.QRCodes,
		},
		Printout:     invoice.PrintoutData,
		Participants: xmlInvoice.Participants(),
	}, nil
}
