package certificates

import (
	"errors"
	"ksef/internal/certsdb"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCertProfile = &cobra.Command{
	Use:   "set-profile",
	Short: "ustaw nazwę profilu dla grupy certyfikatów",
	RunE:  setCertProfileRun,
}

var errEmptyCertIds = errors.New("nie podano identyfikatorów certyfikatów")

const (
	flagNameCertIds     = "cert-ids"
	flagNameProfileName = "profile"
)

func init() {
	flagSet := setCertProfile.Flags()
	flagSet.StringSlice(flagNameCertIds, []string{}, "identyfikatory certyfikatów")
	flagSet.String(flagNameProfileName, "", "nazwa profilu")

	setCertProfile.MarkFlagRequired(flagNameProfileName)
	CertificatesCommand.AddCommand(setCertProfile)
}

func setCertProfileRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	flags := cmd.Flags()
	certIds, err := flags.GetStringSlice(flagNameCertIds)
	if err != nil {
		return err
	}
	if len(certIds) == 0 {
		return errEmptyCertIds
	}
	profileName, err := flags.GetString(flagNameProfileName)
	if err != nil {
		return err
	}

	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}

	if err = certsDB.ForEach(
		func(cert certsdb.Certificate) bool {
			return slices.Contains(certIds, cert.UID)
		},
		func(cert *certsdb.Certificate) error {
			cert.ProfileName = profileName
			return nil
		},
	); err != nil {
		return err
	}

	return certsDB.Save()
}
