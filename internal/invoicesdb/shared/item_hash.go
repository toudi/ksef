package shared

import "strings"

type ItemHash struct {
	Wildcard bool   `yaml:"wildcard,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Index    string `yaml:"index,omitempty"` // Fa -> FaWiersz -> Indeks
	GTIN     string `yaml:"gtin,omitempty"`  // Fa -> FaWiersz -> ..
	PKWiU    string `yaml:"pkwiu,omitempty"` // ...
}

func (h ItemHash) Matches(other ItemHash) bool {
	if h.Wildcard {
		return true
	}

	// let's go from the strongest candidates (i.e. codes) to the weakest ones
	return (compare(h.PKWiU, other.PKWiU) ||
		compare(h.GTIN, other.GTIN) ||
		compare(h.Index, other.Index) ||
		compare(h.Name, other.Name))
}

type JPKItemRule struct {
	Hash         ItemHash `yaml:"hash"`
	Exclude      bool     `yaml:"exclude,omitempty"`
	Vat50Percent bool     `yaml:"vat-50-percent,omitempty"`
	FixedAsset   bool     `yaml:"fixed-asset,omitempty"`
}

func compare(value1, value2 string) bool {
	// reason for the empty string comparison here is that we do not have a guarantee
	// that all of the invoice items will have bunch of these selector fields populated
	// (i.e. PKWiU, GTIN, etc). Therefore if the hash will have this property empty and
	// the comparing invoice item will also have an empty value - that's not good.
	// especially that we're comparing more concrete first before we ever fallback to
	// comparing names themselves which would cause an early return.
	// please refer to lines 19..22 of this file
	return value1 != "" && value2 != "" && strings.EqualFold(value1, value2)
}
