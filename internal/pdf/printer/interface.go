package printer

import "ksef/internal/registry"

type PDFPrinter interface {
	Print(contentBase64 string, meta registry.Invoice, output string) error
	PrintUPO(contentBase64 string, output string) error
}
