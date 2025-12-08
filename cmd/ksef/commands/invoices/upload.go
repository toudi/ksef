package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/internal/invoicesdb"
	invoicesdbconfig "ksef/internal/invoicesdb/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "wyślij faktury do KSeF",
	RunE:  uploadInvoicesRun,
}

func init() {
	invoicesdbconfig.UploaderFlags(uploadCommand.Flags())
	InvoicesCommand.AddCommand(uploadCommand)
}

func uploadInvoicesRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	ksefClient, err := client.InitClient(cmd)
	if err != nil {
		return err
	}
	defer ksefClient.Close()
	// initialize the invoicesdb
	invoicesDB, err := invoicesdb.NewInvoicesDB(vip, invoicesdb.WithKSeFClient(ksefClient))
	if err != nil {
		return err
	}

	return invoicesDB.UploadOutstandingInvoices(cmd.Context(), vip)
}
