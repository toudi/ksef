package invoices

import "github.com/spf13/cobra"

var downloadCommand = &cobra.Command{
	Use:   "download",
	Short: "pobierz faktury",
}

func init() {
	InvoicesCommand.AddCommand(downloadCommand)
}
