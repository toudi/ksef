package certificates

import (
	"errors"
	"ksef/core/certsdb"
	"ksef/internal/resources"
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
	ctx := cmd.Context()

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

	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}

	certs, err := certsDB.FilterCertificates(ctx, func(c certsdb.Certificate) bool {
		return slices.Contains(certIds, c.UID)
	})
	if err != nil {
		return err
	}

	updateData := certsdb.CertificateUpdateData{ProfileName: &profileName}

	for _, cert := range certs {
		if err = certsDB.UpdateCertificate(ctx, cert.UID, updateData); err != nil {
			return err
		}
	}

	return resources.Close(ctx)
}
