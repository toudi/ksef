package jpk

import (
	"fmt"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
	"strconv"

	"github.com/spf13/pflag"
)

const (
	flagNameItems  = "items"
	flagNameGlobal = "global"
)

func itemSelectorFlags(flagSet *pflag.FlagSet) {
	flagSet.StringSlice(flagNameItems, nil, "numery przdmiotów do oznaczenia lub znak *")
	flagSet.BoolP(flagNameGlobal, "g", false, "zapisz dane przedmiotów w ustawieniach podmiotu")
}

func getItemSelector(flagSet *pflag.FlagSet) (config.InvoiceItemSelector, error) {
	var err error
	selector := config.InvoiceItemSelector{}

	if selector.Global, err = flagSet.GetBool(flagNameGlobal); err != nil {
		return selector, err
	}

	if selector.ItemNumbers, err = flagSet.GetStringSlice(flagNameItems); err != nil {
		return selector, err
	}

	return selector, nil
}

func getItemRules(selector config.InvoiceItemSelector, xmlInvoice *monthlyregistry.XMLInvoice, ruleFactory func() shared.JPKItemRule) (rules []shared.JPKItemRule, err error) {
	var itemHash shared.ItemHash

	for _, item := range selector.ItemNumbers {
		if itemHash, err = getItemHash(xmlInvoice, item); err != nil {
			return nil, err
		}
		rule := ruleFactory()
		rule.Hash = itemHash
		rules = append(rules, rule)
	}

	return rules, nil
}

func getItemHash(xmlInvoice *monthlyregistry.XMLInvoice, itemNumber string) (hash shared.ItemHash, err error) {
	if itemNumber == "*" {
		return shared.ItemHash{
			Wildcard: true,
		}, nil
	}

	itemNo, err := strconv.Atoi(itemNumber)
	if err != nil {
		return shared.ItemHash{}, err
	}

	if len(xmlInvoice.Items) < itemNo {
		return shared.ItemHash{}, fmt.Errorf("item %d not found", itemNo)
	}
	itemNo -= 1

	item := xmlInvoice.Items[itemNo]

	if item.GTIN != "" {
		hash.GTIN = item.GTIN
	}
	if item.PKWiU != "" {
		hash.PKWiU = item.PKWiU
	}
	if item.Index != "" {
		hash.Index = item.Index
	}
	hash.Name = item.Name

	return hash, nil
}
