package registry

import (
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry/monthly"
	"ksef/core/invoicedb/invoice"
)

func (fs *FSInvoiceRegistry) GetInvoiceByUID(
	invoiceUID invoicedb.InvoiceUID,
) (invoice *invoicedb.StoredInvoice, err error) {
	var monthlyInvoice *monthly.Invoice
	if invoiceUID.RefNo != nil {
		monthlyInvoice, err = fs.getInvoiceByRefNo(*invoiceUID.RefNo)
	}
	if err != nil {
		return nil, err
	}
	return &invoicedb.StoredInvoice{
		Invoice:   toCoreInvoice(monthlyInvoice),
		YearMonth: fs.yearMonth,
		ID:        monthlyInvoice.ID(),
	}, nil
}

func (fs *FSInvoiceRegistry) getInvoiceByRefNo(refNo string) (*monthly.Invoice, error) {
	monthlyInvoice, err := fs.monthlyRegistry.GetInvoice(
		func(i *monthly.Invoice) bool {
			return i.Type == monthly.InvoiceTypeIssued && i.RefNo == refNo
		},
	)
	if err != nil {
		return nil, err
	}

	return monthlyInvoice, nil
}

func toCoreInvoice(input *monthly.Invoice) *invoice.Invoice {
	return &invoice.Invoice{
		RefNo:     input.RefNo,
		KSeFRefNo: input.KSeFRefNo,
		QRCodes:   invoice.InvoiceQRCodes(input.QRCodes),
		Type:      invoice.InvoiceType(input.Type),
	}
}
