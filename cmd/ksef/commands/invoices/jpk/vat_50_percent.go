package jpk

import (
	"ksef/internal/invoicesdb/shared"

	"github.com/spf13/cobra"
)

var jpk50PercVAT = &cobra.Command{
	Use:     "vat-50 [faktura.xml]",
	Short:   "oznacza pozycje na fakturze zakupowej w JPK jako rozliczane w trybie 50% VAT",
	RunE:    markVat50Percent,
	Args:    cobra.ExactArgs(1),
	PreRunE: initJPKManagerFromInvoiceFile,
}

func init() {
	flagSet := jpk50PercVAT.Flags()
	itemSelectorFlags(flagSet)
}

func markVat50Percent(cmd *cobra.Command, args []string) error {
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
