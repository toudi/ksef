package invoices

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb"
	"ksef/internal/invoicesdb/config"
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
	var flagSet = syncCommand.Flags()
	flags.NIP(flagSet)
	config.SyncFlags(flagSet)
	syncCommand.MarkFlagRequired(flags.FlagNameNIP)
	InvoicesCommand.AddCommand(syncCommand)
}

func syncInvoicesRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}
	invoicesDB, err := invoicesdb.OpenForNIP(nip, vip)
	if err != nil {
		return err
	}
	return invoicesDB.Sync(
		cmd.Context(),
		vip,
	)
}
