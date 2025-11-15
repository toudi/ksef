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
	var env = runtime.GetGateway(viper.GetViper())

	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return err
	}

	httpClient := http.NewClient(string(runtime.GetGateway(viper.GetViper())))

	certificates, err := security.DownloadCertificates(cmd.Context(), httpClient)
	if err != nil {
		return err
	}

	for _, cert := range certificates {
		certHash := certsdb.CertificateHash{
			Usage:       cert.Usage,
			Environment: env,
			ValidFrom:   cert.ValidFrom,
			ValidTo:     cert.ValidTo,
		}

		if err = certsDB.AddIfHashNotFound(certHash, func(newCert *certsdb.Certificate) error {
			logging.CertsDBLogger.Debug("zapisywanie certyfikatu ministerstwa finansów", "usage", certHash.UsageAsString())
			var err error
			var derBytes []byte
			if derBytes, err = base64.StdEncoding.DecodeString(cert.Certificate); err != nil {
				return err
			}
			newCert.Environment = env
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
