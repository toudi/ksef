package certificates

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/logging"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagAuth    = "auth"
	flagOffline = "offline"
)

var prepareCSRCommand = &cobra.Command{
	Use:   "prepare-csr",
	Short: "Przygotuj wnioski certyfikacyjne",
	RunE:  sendCsrs,
}

func init() {
	flagSet := prepareCSRCommand.Flags()
	flags.AuthMethod(prepareCSRCommand)
	flagSet.BoolP(flagAuth, "a", false, "przygotuj wniosek dla certyfikatu autoryzacyjnego")
	flagSet.BoolP(flagOffline, "o", false, "przygotuj wniosek dla certyfikatu offline")

	flags.NIP(prepareCSRCommand.Flags())

	flagSet.SortFlags = false
	prepareCSRCommand.MarkFlagRequired(flags.FlagNameNIP)
	prepareCSRCommand.MarkFlagsOneRequired(flagAuth, flagOffline)
	CertificatesCommand.AddCommand(prepareCSRCommand)
}

func sendCsrs(cmd *cobra.Command, _ []string) error {
	env := config.GetGateway(viper.GetViper())
	nip, _ := cmd.Flags().GetString(flags.FlagNameNIP)
	if cli, err = client.InitClient(cmd); err != nil {
		return err
	}
	certsManager, err := cli.Certificates(env)
	if err != nil {
		return err
	}
	defer certsManager.SaveDB()
	ed, err := certsManager.GetEnrollmentsData(cmd.Context())
	if err != nil {
		return err
	}
	if prepareAuth, _ := cmd.Flags().GetBool(flagAuth); prepareAuth {
		logging.CertsDBLogger.Debug("przygotowuję CSR dla certyfikatu autoryzacji")
		if err = certsManager.PrepareEnrollmentCSR(ed, certsdb.UsageAuthentication, nip); err != nil {
			return err
		}
	}
	if prepareOffline, _ := cmd.Flags().GetBool(flagOffline); prepareOffline {
		logging.CertsDBLogger.Debug("przygotowuję CSR dla certyfikatu offline")
		if err = certsManager.PrepareEnrollmentCSR(ed, certsdb.UsageOffline, nip); err != nil {
			return err
		}
	}
	return nil
	// return cli.HandleEnrollmentsData(cmd.Context())
}
