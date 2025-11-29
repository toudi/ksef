package uploader

func (idb *InvoiceDB) GetByRefNo(refNo string) *Invoice {
	if index, exists := idb.invoiceByRefNoIndex[refNo]; exists {
		return idb.Invoices[index]
	}
	return nil
}
