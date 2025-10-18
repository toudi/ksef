package certificates

import (
	"encoding/base64"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/environment"
	"ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client/v2/security"

	"github.com/spf13/cobra"
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
	env := environment.FromContext(cmd.Context())
	cfg := config.GetConfig().APIConfig(env)
	certsDB := cfg.CertificatesDB

	httpClient := http.NewClient(cfg.Environment.Host)
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
