package printer

import monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

type PDFPrinter interface {
	PrintInvoice(
		srcFile string, output string,
		meta *monthlyregistry.InvoicePrintingMeta,
	) error
	PrintUPO(srcFile string, output string) error
}
