package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/logging"
	"ksef/internal/resources"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "synchronizuj faktury z KSeF",
	RunE:  syncInvoicesRun,
}

func init() {
	flagSet := syncCommand.Flags()
	flags.NIP(flagSet)
	config.SyncFlags(flagSet)
	runtime.CertProfileFlag(flagSet)
	InvoicesCommand.AddCommand(syncCommand)
}

func syncInvoicesRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	vip := viper.GetViper()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	nip, err := runtime.GetNIP(vip)
	if err != nil {
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

	invoicesDB, err := invoicesdb.OpenForNIP(
		nip,
		vip,
		certsDB,
		invoicesdb.WithKSeFClient(ksefClient),
		invoicesdb.WithMonthsRangeGenerator(invoicesdb.MonthsRangePreviousMonthGenerator),
	)
	if err != nil {
		return err
	}
	if err = invoicesDB.Sync(
		ctx,
		vip,
	); err != nil {
		return err
	}

	return ksefClient.WaitForTokenManagerLoop()
}
