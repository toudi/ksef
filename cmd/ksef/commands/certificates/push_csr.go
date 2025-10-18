package certificates

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/environment"

	"github.com/spf13/cobra"
)

var syncEnrollmentsCommand = &cobra.Command{
	Use:   "sync-csr",
	Short: "wysyła przygotowane wnioski CSR oraz pobiera gotowe certyfikaty",
	RunE:  syncEnrollments,
}

func init() {
	flags.AuthMethod(syncEnrollmentsCommand)
	flags.NIP(syncEnrollmentsCommand)
	syncEnrollmentsCommand.Flags().SortFlags = false
	CertificatesCommand.AddCommand(syncEnrollmentsCommand)
}

func syncEnrollments(cmd *cobra.Command, _ []string) error {
	if cli, err = client.InitClient(cmd); err != nil {
		return err
	}
	certsManager, err := cli.Certificates(environment.FromContext(cmd.Context()))
	if err != nil {
		return err
	}
	defer certsManager.SaveDB()
	return certsManager.SyncEnrollments(cmd.Context())
}
