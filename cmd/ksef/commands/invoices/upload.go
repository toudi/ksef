package invoices

import "github.com/spf13/cobra"

var uploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "wyślij faktury do KSeF",
}

func init() {
	InvoicesCommand.AddCommand(uploadCommand)
}
