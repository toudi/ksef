package certificates

import (
	"encoding/base64"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/security"
	"ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadKSeFCertsCommand = &cobra.Command{
	Use:   "download",
	Short: "Pobranie certyfikatów klucza publicznego MF",
	RunE:  downloadKSeFCerts,
}

func init() {
	CertificatesCommand.AddCommand(downloadKSeFCertsCommand)
}

func downloadKSeFCerts(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	env := runtime.GetEnvironment(vip)

	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}

	httpClient := http.NewClient(env.API)

	certificates, err := security.DownloadCertificates(cmd.Context(), httpClient)
	if err != nil {
		return err
	}

	for _, cert := range certificates {
		certHash := certsdb.CertificateHash{
			Usage:         cert.Usage,
			EnvironmentId: env.ID,
			ValidFrom:     cert.ValidFrom,
			ValidTo:       cert.ValidTo,
		}

		if err = certsDB.AddIfHashNotFound(certHash, func(newCert *certsdb.Certificate) error {
			logging.CertsDBLogger.Debug("zapisywanie certyfikatu ministerstwa finansów", "usage", certHash.UsageAsString())
			var err error
			var derBytes []byte
			if derBytes, err = base64.StdEncoding.DecodeString(cert.Certificate); err != nil {
				return err
			}
			newCert.EnvironmentId = env.ID
			newCert.CertificateHash = certHash
			if err = newCert.SaveCert(derBytes); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return certsDB.Save()
}
