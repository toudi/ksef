package monthlyregistry

func (i *Invoice) GetPrintingMeta() *InvoicePrintingMeta {
	return &InvoicePrintingMeta{
		Usage:    "invoice:" + invoiceTypeToPrinterUsage[i.Type],
		Printout: i.PrintoutData,
		Invoice: InvoiceMeta{
			KSeFRefNo: i.KSeFRefNo,
			QRCodes:   i.QRCodes,
		},
	}
}
