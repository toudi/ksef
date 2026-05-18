package shared

import (
	"fmt"
	"ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
	"strconv"

	"github.com/spf13/pflag"
)

const (
	FlagNameItems  = "items"
	FlagNameGlobal = "global"
)

func ItemSelectorFlags(flagSet *pflag.FlagSet) {
	flagSet.StringSlice(FlagNameItems, nil, "numery przedmiotów do oznaczenia lub znak *")
	flagSet.BoolP(FlagNameGlobal, "g", false, "zapisz dane przedmiotów w ustawieniach podmiotu")
}

func GetItemSelector(flagSet *pflag.FlagSet) (config.InvoiceItemSelector, error) {
	var err error
	selector := config.InvoiceItemSelector{}

	if selector.Global, err = flagSet.GetBool(FlagNameGlobal); err != nil {
		return selector, err
	}

	if selector.ItemNumbers, err = flagSet.GetStringSlice(FlagNameItems); err != nil {
		return selector, err
	}

	return selector, nil
}

func GetItemHash(xmlInvoice *monthlyregistry.XMLInvoice, itemNumber string) (shared.ItemHash, error) {
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
		return shared.ItemHash{GTIN: item.GTIN}, nil
	}
	if item.PKWiU != "" {
		return shared.ItemHash{PKWiU: item.PKWiU}, nil
	}
	if item.Index != "" {
		return shared.ItemHash{Index: item.Index}, nil
	}
	return shared.ItemHash{Name: item.Name}, nil
}
