package types

// defining this structure just in case other schemas are even worse to fill out.
// for instance - right now there's a field P_14 that seems to have something to do
// with railway which we don't support anyway. However, who knows what the ministry
// will think of in the next iterations. so basically just to be safe I'd reserve
// this struct so that we could extend it with additional properties
type SaleAttributes struct {
	VatRate VatRate
}

type Sales struct {
	ByAttributes map[SaleAttributes]*VATInfo
	Total        *VATInfo
}

func (s *Sales) Add(attributes SaleAttributes, info VATInfo) {
	if s.ByAttributes == nil {
		s.ByAttributes = make(map[SaleAttributes]*VATInfo)
		s.Total = &VATInfo{}
	}
	if _, exists := s.ByAttributes[attributes]; !exists {
		s.ByAttributes[attributes] = &VATInfo{}
	}
	s.ByAttributes[attributes].Add(info)
	s.Total.Add(info)
}
