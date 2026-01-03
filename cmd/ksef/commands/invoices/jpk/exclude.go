package jpk

import (
	"ksef/internal/invoicesdb/shared"

	"github.com/spf13/cobra"
)

var jpkExclude = &cobra.Command{
	Use:     "exclude [faktura.xml]",
	Short:   "oznacza pozycje faktury zakupowej jako wyłączone z JPK",
	RunE:    markInvoiceExcluded,
	Args:    cobra.ExactArgs(1),
	PreRunE: initJPKManagerFromInvoiceFile,
}

func init() {
	flagSet := jpkExclude.Flags()
	itemSelectorFlags(flagSet)
}

func markInvoiceExcluded(cmd *cobra.Command, args []string) error {
	itemSelector, err := getItemSelector(cmd.Flags())
	if err != nil {
		return err
	}

	rules, err := getItemRules(
		itemSelector,
		xmlInvoice, func() shared.JPKItemRule {
			return shared.JPKItemRule{Exclude: true}
		},
	)
	if err != nil {
		return err
	}

	return jpkManager.AddItemRules(
		invoice, rules, itemSelector.Global,
	)
}
