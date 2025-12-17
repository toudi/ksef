package printer

import (
	"ksef/internal/registry/types"
)

type PDFPrinter interface {
	Print(contentBase64 string, meta types.Invoice, output string) error
	PrintUPO(srcFile string, output string) error
}
