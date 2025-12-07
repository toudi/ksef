package invoices

import "github.com/spf13/cobra"

var upoCommand = &cobra.Command{
	Use:   "upo",
	Short: "pobierz UPO dla wysłanych faktur",
}

func init() {
	InvoicesCommand.AddCommand(upoCommand)
}
