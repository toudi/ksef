package uploader

import monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

func (u *Uploader) AddToQueue(invoiceWithMeta *monthlyregistry.InvoiceMetadata) {
	if _, exists := u.Queue[invoiceWithMeta.FormCode]; !exists {
		u.Queue[invoiceWithMeta.FormCode] = make([]string, 0)
	}
	u.Queue[invoiceWithMeta.FormCode] = append(u.Queue[invoiceWithMeta.FormCode], invoiceWithMeta.Filename)
	u.registryByInvoiceChecksum[invoiceWithMeta.Invoice.Checksum] = invoiceWithMeta.Registry
}
