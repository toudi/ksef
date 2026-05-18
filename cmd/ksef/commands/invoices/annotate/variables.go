package annotate

import (
	sharedcli "ksef/cmd/ksef/commands/invoices/shared"
	"ksef/internal/invoicesdb/shared"

	"github.com/spf13/cobra"
)

var (
	ctx *sharedcli.InvoiceContext
)

func initAnnotationRule(cmd *cobra.Command, args []string) error {
	var err error
	ctx, err = sharedcli.InitAnnotationsManagerFromInvoiceFile(cmd, args)
	return err
}

func buildAnnotationRule(cmd *cobra.Command, flagValue string) ([]shared.Annotation, error) {
	selector, err := sharedcli.GetItemSelector(cmd.Flags())
	if err != nil {
		return nil, err
	}

	var rules []shared.Annotation
	for _, item := range selector.ItemNumbers {
		hash, err := sharedcli.GetItemHash(ctx.XMLInvoice, item)
		if err != nil {
			return nil, err
		}
		rules = append(rules, shared.Annotation{
			Hash:    hash,
			Comment: &flagValue,
		})
	}

	return rules, nil
}
