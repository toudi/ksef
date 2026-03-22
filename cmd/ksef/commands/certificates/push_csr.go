package certificates

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	kr "ksef/internal/keyring"
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
	keyring, err := kr.NewKeyring(viper.GetViper())
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	defer keyring.Close()

	if cli, err = client.InitClient(cmd, viper.GetViper(), keyring); err != nil {
		return err
	}
	certsManager, err := cli.Certificates(runtime.GetEnvironmentId(viper.GetViper()))
	if err != nil {
		return err
	}
	defer certsManager.SaveDB()
	return certsManager.SyncEnrollments(cmd.Context())
}
