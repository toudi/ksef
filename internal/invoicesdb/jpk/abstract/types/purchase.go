package types

import (
	"ksef/internal/invoicesdb/jpk/abstract/processors/vat"
	"ksef/internal/invoicesdb/jpk/manager"
	"ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

type SellerData struct {
	NIP  string
	Name string
}

type PurchaseItem struct {
	purchase   *Purchase
	RefNo      string
	KSeFRefNo  string
	Date       string
	Seller     SellerData
	VATAmounts *types.PurchaseVAT
}

type Purchase struct {
	Rows       []*PurchaseItem
	processors []PurchaseInvoiceProcessor
	VATAmounts *types.PurchaseVAT
	manager    *manager.JPKManager
}

func NewPurchase(manager *manager.JPKManager, processors []PurchaseInvoiceProcessor) *Purchase {
	return &Purchase{
		processors: processors,
		manager:    manager,
		VATAmounts: &types.PurchaseVAT{},
	}
}

func (p *Purchase) ProcessInvoice(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
) error {
	row := &PurchaseItem{
		purchase:   p,
		VATAmounts: &types.PurchaseVAT{},
	}

	for _, processor := range p.processors {
		if err := processor(p.manager, invoice, doc, row); err != nil {
			return err
		}
	}

	p.Rows = append(p.Rows, row)

	return nil
}

func (p *PurchaseItem) AddAmount(attributes types.PurchaseAttributes, vatInfo vat.VatInfo) {
	vatAmounts := types.VATInfo{
		Base: vatInfo.NetAmount,
		Vat:  vatInfo.VatAmount,
	}
	p.VATAmounts.Add(attributes, vatAmounts)
	p.purchase.VATAmounts.Add(attributes, vatAmounts)
}
