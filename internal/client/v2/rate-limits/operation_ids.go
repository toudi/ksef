package ratelimits

// no idea why the ministry couldn't have simply used HTTP paths - then we could
// sort them by length and use the best fit
const (
	OperationInvoicesExport  = "invoiceExport"
	OperationExportStatus    = "invoiceExportStatus"
	OperationInvoiceSend     = "invoiceSend"
	OperationInvoiceMetadata = "invoiceMetadata"
	OperationInvoiceDownload = "invoiceDownload"
)
