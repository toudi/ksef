package registry

import (
	"errors"
	"fmt"
)

var ErrPaymentIdNotFound = errors.New("unable to find payment with the specified ID")

func (r *InvoiceRegistry) getInvoiceByRefNo(invoiceId string) (*Invoice, error) {
	var index int
	var exists bool
	index, exists = r.seiRefNoIndex[invoiceId]
	if exists {
		return &r.Invoices[index], nil
	}
	index, exists = r.refNoIndex[invoiceId]
	if exists {
		return &r.Invoices[index], nil
	}

	return nil, fmt.Errorf("unable to find invoice")
}

func (r *InvoiceRegistry) GetSEIRefNoFromArray(invoiceIds []string) ([]string, error) {
	result := make([]string, len(invoiceIds))

	for index, invoiceId := range invoiceIds {
		invoice, err := r.getInvoiceByRefNo(invoiceId)
		if err != nil {
			return nil, fmt.Errorf("unable to find invoice with the following ID in the registry: %v", invoiceId)
		}

		result[index] = invoice.KSeFReferenceNumber
	}

	return result, nil
}

func (r *InvoiceRegistry) GetInvoiceIdsForPaymentId(paymentId string) ([]InvoiceRefId, error) {
	result := make([]InvoiceRefId, 0)

	for _, payment := range r.PaymentIds {
		if payment.SEIPaymentRefNo == paymentId {
			for _, invoiceId := range payment.InvoiceIDS {
				invoice, err := r.getInvoiceByRefNo(invoiceId)
				if err != nil {
					return nil, fmt.Errorf("unable to find invoice with the specified ID: %s", invoiceId)
				}
				result = append(result, InvoiceRefId{
					ReferenceNumber:     invoice.ReferenceNumber,
					KSeFReferenceNumber: invoice.KSeFReferenceNumber,
				})
			}

			return result, nil
		}
	}

	return nil, ErrPaymentIdNotFound
}
