package invoices

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	"ksef/internal/invoicesdb/config"
	kr "ksef/internal/keyring"
	"ksef/internal/logging"
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
	vip := viper.GetViper()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	keyring, err := kr.NewKeyring(vip)
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	defer keyring.Close()

	ksefClient, err := client.InitClient(cmd, vip, keyring)
	if err != nil {
		return err
	}

	invoicesDB, err := invoicesdb.OpenForNIP(nip, vip, invoicesdb.WithKSeFClient(ksefClient))
	if err != nil {
		return err
	}
	if err = invoicesDB.Sync(
		cmd.Context(),
		vip,
	); err != nil {
		return err
	}

	return ksefClient.WaitForTokenManagerLoop()
}
