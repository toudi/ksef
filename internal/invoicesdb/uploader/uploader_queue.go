package uploader

import (
	sessionTypes "ksef/internal/client/v2/session/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

func (u *Uploader) AddToQueue(invoiceWithMeta *monthlyregistry.InvoiceMetadata) {
	if _, exists := u.Queue[invoiceWithMeta.FormCode]; !exists {
		u.Queue[invoiceWithMeta.FormCode] = make([]sessionTypes.Invoice, 0)
	}
	u.Queue[invoiceWithMeta.FormCode] = append(
		u.Queue[invoiceWithMeta.FormCode],
		sessionTypes.Invoice{
			Checksum: invoiceWithMeta.Invoice.Checksum,
			Offline:  invoiceWithMeta.Invoice.Offline,
			Filename: invoiceWithMeta.Filename,
		},
	)
}
