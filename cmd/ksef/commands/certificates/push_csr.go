package certificates

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	kr "ksef/core/keyring"
	"ksef/internal/logging"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncEnrollmentsCommand = &cobra.Command{
	Use:   "sync-csr",
	Short: "wysyła przygotowane wnioski CSR oraz pobiera gotowe certyfikaty",
	RunE:  syncEnrollments,
}

func init() {
	flags.NIP(syncEnrollmentsCommand.Flags())
	syncEnrollmentsCommand.Flags().SortFlags = false
	syncEnrollmentsCommand.MarkFlagRequired(flags.FlagNameNIP)
	CertificatesCommand.AddCommand(syncEnrollmentsCommand)
}

func syncEnrollments(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	keyring, err := kr.NewKeyring(vip)
	ctx := cmd.Context()
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	defer keyring.Close()

	if cli, err = client.InitClient(cmd, vip, keyring); err != nil {
		return err
	}
	certsManager, err := cli.Certificates(runtime.GetEnvironmentId(vip))
	if err != nil {
		return err
	}
	defer certsManager.SaveDB(ctx)
	return certsManager.SyncEnrollments(ctx)
}
