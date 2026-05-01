package abstract

import (
	"ksef/internal/invoicesdb/jpk/abstract/processors/purchase"
	"ksef/internal/invoicesdb/jpk/abstract/processors/sale"
	"ksef/internal/invoicesdb/jpk/abstract/types"
	"ksef/internal/invoicesdb/jpk/manager"
)

// this module contains abstract information about Sale and Purchase rows.
// it will be consumed by dedicated generators and converted into appropriate fields.
// abstract code will therefore **NOT** contain any form-specific fields.
type MonthlyReport struct {
	Sales    *types.Sales
	Purchase *types.Purchase
}

var saleInvoiceProcessors = []types.SaleInvoiceProcessor{
	sale.ExtractRefNos,
	sale.ExtractBuyer,
	sale.ExtractDates,
	sale.AggregateAmountsByVATRate,
}

var purchaseInvoiceProcessors = []types.PurchaseInvoiceProcessor{
	purchase.ExtractRefNos,
	purchase.ExtractSeller,
	purchase.AggregateItems,
}

func NewMonthlyReport(manager *manager.JPKManager) *MonthlyReport {
	return &MonthlyReport{
		Sales:    types.NewSales(saleInvoiceProcessors),
		Purchase: types.NewPurchase(manager, purchaseInvoiceProcessors),
	}
}
