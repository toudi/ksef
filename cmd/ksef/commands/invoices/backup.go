package invoices

import "github.com/spf13/cobra"

var backupCommand = &cobra.Command{
	Use:   "backup",
	Short: "zarchiwizuj bazę faktur",
}

func init() {
	InvoicesCommand.AddCommand(backupCommand)
}
