package invoices

import (
	"ksef/cmd/ksef/commands/invoices/jpk"
	"ksef/internal/invoicesdb/config"

	"github.com/spf13/cobra"
)

var InvoicesCommand = &cobra.Command{
	Use:   "invoices",
	Short: "zarządzanie bazą faktur",
}

const (
	flagNameConfirm = "confirm"
)

func init() {
	config.InvoicesDBFlags(InvoicesCommand.PersistentFlags())
	InvoicesCommand.PersistentFlags().Bool(flagNameConfirm, false, "potwierdź operację")
	InvoicesCommand.AddCommand(jpk.JPKCommand)
}
