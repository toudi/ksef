package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"ksef/internal/logging"
	"ksef/internal/resources"
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
	runtime.CertProfileFlag(uploadCommand.Flags())
	InvoicesCommand.AddCommand(uploadCommand)
}

func uploadInvoicesRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	ctx := cmd.Context()

	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	keyring, err := resources.Keyring(vip)
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)

	ksefClient, err := client.InitClient(cmd, vip, keyring)
	if err != nil {
		return err
	}
	defer ksefClient.Close()

	// initialize the invoicesdb
	invoicesDB, err := invoicesdb.NewInvoicesDB(
		vip,
		certsDB,
		invoicesdb.WithKSeFClient(ksefClient),
		invoicesdb.WithKeyring(keyring),
		invoicesdb.WithMonthsRangeGenerator(invoicesdb.MonthsRangePreviousMonthGenerator),
	)
	if err != nil {
		return err
	}

	return invoicesDB.UploadOutstandingInvoices(
		ctx,
		uploaderconfig.GetUploaderConfig(vip),
		statuscheckerconfig.GetStatusCheckerConfig(vip),
	)
}
