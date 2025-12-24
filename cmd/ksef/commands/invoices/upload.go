package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "wyślij faktury do KSeF",
	RunE:  uploadInvoicesRun,
	Args:  cobra.MaximumNArgs(1),
}

func init() {
	uploaderconfig.UploaderFlags(uploadCommand.Flags())
	statuscheckerconfig.StatusCheckerFlags(uploadCommand.Flags())
	flags.NIP(uploadCommand.Flags())
	uploadCommand.MarkFlagRequired(flags.FlagNameNIP)
	runtime.CertProfileFlag(uploadCommand.Flags())
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

	return invoicesDB.UploadOutstandingInvoices(
		cmd.Context(),
		uploaderconfig.GetUploaderConfig(vip),
		statuscheckerconfig.GetStatusCheckerConfig(vip),
	)
}
