package registry

import (
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/fs/metadata/registry/monthly"
)

func (fs *FSInvoiceRegistry) GetInvoiceByUID(
	invoiceUID invoicedb.InvoiceUID,
) (*invoicedb.StoredInvoice, error) {
	var monthlyInvoice *monthly.Invoice
	var err error

	// TODO: add future support for retrieving invoice via KSeF Ref no
	// that's why the code looks rather awkward right now
	if invoiceUID.RefNo != nil {
		monthlyInvoice, err = fs.getInvoiceByRefNo(*invoiceUID.RefNo)
	} else {
		return nil, invoicedb.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if monthlyInvoice == nil {
		return nil, invoicedb.ErrNotFound
	}

	return &invoicedb.StoredInvoice{
		NIP:       fs.nip,
		Invoice:   monthly.ToCoreInvoice(monthlyInvoice),
		YearMonth: fs.yearMonth,
		ID:        monthlyInvoice.ID(),
	}, nil
}

func (fs *FSInvoiceRegistry) getInvoiceByRefNo(refNo string) (*monthly.Invoice, error) {
	return fs.monthlyRegistry.GetInvoice(
		func(i *monthly.Invoice) bool {
			return i.Type == monthly.InvoiceTypeIssued && i.RefNo == refNo
		},
	)
}
