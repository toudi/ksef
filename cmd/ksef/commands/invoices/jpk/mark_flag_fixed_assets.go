package jpk

import (
	"ksef/internal/invoicesdb/shared"

	"github.com/spf13/cobra"
)

var jpkFixedAssets = &cobra.Command{
	Use:     "fixed-assets [faktura.xml]",
	Short:   "oznacza pozycje faktury zakupowej jako środki trwałe w JPK",
	RunE:    markFixedAssets,
	Args:    cobra.ExactArgs(1),
	PreRunE: initJPKManagerFromInvoiceFile,
}

func init() {
	flagSet := jpkFixedAssets.Flags()
	itemSelectorFlags(flagSet)
}

func markFixedAssets(cmd *cobra.Command, args []string) error {
	itemSelector, err := getItemSelector(cmd.Flags())
	if err != nil {
		return err
	}

	rules, err := getItemRules(
		itemSelector,
		xmlInvoice, func() shared.JPKItemRule {
			return shared.JPKItemRule{FixedAsset: true}
		},
	)
	if err != nil {
		return err
	}

	return jpkManager.AddItemRules(
		invoice, rules, itemSelector.Global,
	)
}
