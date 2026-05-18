package annotate

import (
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:     "clear [faktura.xml]",
	Short:   "usuwa wszystkie adnotacje z faktury",
	RunE:    runClear,
	Args:    cobra.ExactArgs(1),
	PreRunE: initAnnotationRule,
}

func runClear(cmd *cobra.Command, args []string) error {
	ctx.Invoice.Annotations = nil
	return ctx.InvoiceRegistry.Save()
}
