package shared

import (
	"regexp"
	"strings"
)

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
		compareName(h.Name, other.Name))
}

type Annotation struct {
	Hash         ItemHash `yaml:"hash"`
	Exclude      bool     `yaml:"exclude,omitempty"`
	Vat50Percent bool     `yaml:"vat-50-percent,omitempty"`
	FixedAsset   bool     `yaml:"fixed-asset,omitempty"`
	Comment      *string  `yaml:"comment,omitempty"`
}

func (a Annotation) String() string {
	var flags []string
	if a.Exclude {
		flags = append(flags, "wyłącz z raportu (zakup prywatny)")
	}
	if a.Vat50Percent {
		flags = append(flags, "50% VAT")
	}
	if a.FixedAsset {
		flags = append(flags, "środki trwałe")
	}
	if a.Comment != nil {
		flags = append(flags, *a.Comment)
	}
	return strings.Join(flags, ", ")
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

func compareName(value1, value2 string) bool {
	// the invoice items can be wildcarded so let's use regexp to compare them
	match, _ := regexp.MatchString(value1, value2)
	if match {
		return true
	}
	match, _ = regexp.MatchString(value2, value1)
	return match
}
