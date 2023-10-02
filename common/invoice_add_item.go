package common

import "fmt"

func (i *Invoice) AddItem(item *InvoiceItem) error {
	if item.UnitPrice.Vat.Description == "" {
		item.UnitPrice.Vat.Description = fmt.Sprintf("%d", item.UnitPrice.Vat.Rate)
	}
	i.Items = append(i.Items, item)

	// increment total amounts.
	if i.TotalPerVATRate == nil {
		i.TotalPerVATRate = make(map[string]Amount)
	}

	amount := item.Amount()

	var total Amount

	total, exists := i.TotalPerVATRate[item.UnitPrice.Vat.Description]
	if !exists {
		total = Amount{}
	}

	total.Net += amount.Net
	total.Gross += amount.Gross
	total.VAT += amount.VAT

	i.TotalPerVATRate[item.UnitPrice.Vat.Description] = total

	i.Total.Net += amount.Net
	i.Total.Gross += amount.Gross
	i.Total.VAT += amount.VAT

	return nil
}
