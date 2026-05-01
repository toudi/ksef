package types

type PurchaseAttributes struct {
	VatRate     VatRate
	FixedAssets bool
}

type PurchaseVAT struct {
	ByAttributes map[PurchaseAttributes]*VATInfo
	Total        *VATInfo
}

func (p *PurchaseVAT) Add(attributes PurchaseAttributes, info VATInfo) {
	if p.ByAttributes == nil {
		p.ByAttributes = make(map[PurchaseAttributes]*VATInfo)
		p.Total = &VATInfo{}
	}
	if _, exists := p.ByAttributes[attributes]; !exists {
		p.ByAttributes[attributes] = &VATInfo{}
	}
	p.ByAttributes[attributes].Add(info)
	p.Total.Add(info)
}
