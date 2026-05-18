package annotate

import (
	sharedcli "ksef/cmd/ksef/commands/invoices/shared"
	"ksef/internal/invoicesdb/shared"

	"github.com/spf13/cobra"
)

var annotateExcludeCmd = &cobra.Command{
	Use:     "exclude [faktura.xml]",
	Short:   "oznacza pozycje faktury zakupowej jako wyłączone z JPK",
	RunE:    runExclude,
	Args:    cobra.ExactArgs(1),
	PreRunE: initAnnotationRule,
}

func init() {
	flagSet := annotateExcludeCmd.Flags()
	sharedcli.ItemSelectorFlags(flagSet)
}

func runExclude(cmd *cobra.Command, args []string) error {
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
			Hash:    hash,
			Exclude: true,
		})
	}

	return ctx.AnnotationsMgr.AddItemRules(ctx.Invoice, rules, global)
}
