package sessionregistry

import (
	"ksef/internal/client/v2/session/status"
	"ksef/internal/utils"

	"github.com/samber/lo"
)

func (us *UploadSession) addInfoAboutInvoice(
	uploadStatus status.InvoiceInfo,
) { // yeah, yeah. I should implement support for indices. I get it.
	invoice, found := lo.Find(us.Invoices, func(invoice *Invoice) bool {
		return invoice.Checksum == uploadStatus.Checksum
	})
	if found {
		invoice.RefNo = uploadStatus.InvoiceNumber
		invoice.KSeFRefNo = uploadStatus.KSeFRefNo
		if len(uploadStatus.Status.Details) > 0 {
			invoice.Errors = uploadStatus.Status.Details
		}
	} else {
		var invoice = &Invoice{
			RefNo:    uploadStatus.InvoiceNumber,
			Checksum: utils.Base64ToHex(uploadStatus.Checksum),
		}
		if uploadStatus.Status.Successful() {
			invoice.KSeFRefNo = uploadStatus.KSeFRefNo
		} else {
			invoice.Errors = uploadStatus.Status.Details
		}
		us.Invoices = append(us.Invoices, invoice)
	}
}
