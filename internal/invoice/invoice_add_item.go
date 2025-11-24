package invoice

import "fmt"

func (i *Invoice) AddItem(item *InvoiceItem) error {
	if item.UnitPrice.Vat.Description == "" {
		item.UnitPrice.Vat.Description = fmt.Sprintf("%d", item.UnitPrice.Vat.Rate)
	}
	i.Items = append(i.Items, item)

	amount := item.Amount()

	i.adjustTotalAmount(amount, item.UnitPrice.Vat.Description)

	return nil
}

func (i *Invoice) adjustTotalAmount(amount Amount, vatDescription string) {
	// increment total amounts.
	if i.TotalPerVATRate == nil {
		i.TotalPerVATRate = make(map[string]Amount)
	}

	var total Amount

	total, exists := i.TotalPerVATRate[vatDescription]
	if !exists {
		total = Amount{}
	}

	total.Net += amount.Net
	total.Gross += amount.Gross
	total.VAT += amount.VAT

	i.TotalPerVATRate[vatDescription] = total

	i.Total.Net += amount.Net
	i.Total.Gross += amount.Gross
	i.Total.VAT += amount.VAT

}

func (i *Invoice) AddCorrectedItem(oldItem, newItem *InvoiceItem) error {
	oldAmount := oldItem.Amount()
	oldAmount.Net *= -1
	oldAmount.Gross *= -1
	oldAmount.VAT *= -1
	newAmount := newItem.Amount()

	i.Items = append(i.Items, oldItem, newItem)

	i.adjustTotalAmount(oldAmount, oldItem.UnitPrice.Vat.Description)
	i.adjustTotalAmount(newAmount, newItem.UnitPrice.Vat.Description)

	return nil
}
