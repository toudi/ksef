package annotate

import (
	sharedcli "ksef/cmd/ksef/commands/invoices/shared"
	"ksef/internal/invoicesdb/shared"

	"github.com/spf13/cobra"
)

var annotateVat50Cmd = &cobra.Command{
	Use:     "vat-50 [faktura.xml]",
	Short:   "oznacza pozycje faktury zakupowej jako podlegające 50% VAT",
	RunE:    runVat50,
	Args:    cobra.ExactArgs(1),
	PreRunE: initAnnotationRule,
}

func init() {
	flagSet := annotateVat50Cmd.Flags()
	sharedcli.ItemSelectorFlags(flagSet)
}

func runVat50(cmd *cobra.Command, args []string) error {
	selector, err := sharedcli.GetItemSelector(cmd.Flags())
	if err != nil {
		return err
	}

	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	var rules []shared.Annotation
	for _, item := range selector.ItemNumbers {
		hash, err := sharedcli.GetItemHash(ctx.XMLInvoice, item)
		if err != nil {
			return err
		}
		rules = append(rules, shared.Annotation{
			Hash:         hash,
			Vat50Percent: true,
		})
	}

	return ctx.AnnotationsMgr.AddItemRules(ctx.Invoice, rules, global)
}
