package jpk

import (
	"github.com/spf13/cobra"
)

var jpkClearFlags = &cobra.Command{
	Use:     "clear [faktura.xml]",
	Short:   "usuwa flagi JPK z faktury zakupowej",
	RunE:    clearJPKFlags,
	Args:    cobra.ExactArgs(1),
	PreRunE: initJPKManagerFromInvoiceFile,
}

func clearJPKFlags(cmd *cobra.Command, args []string) error {
	invoice.JPK.ItemRules = nil
	return invoiceRegistry.Save()
}
