package certificates

import (
	"encoding/base64"
	"ksef/cmd/ksef/commands/client"
	"ksef/core/certsdb"
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
	ctx := cmd.Context()

	gatewayConfig, err := runtime.GatewayConfig(vip)
	if err != nil {
		return err
	}

	certsDB, err := client.CertsDB(vip)
	if err != nil {
		return err
	}

	httpClient := http.NewClient(gatewayConfig.API)

	certificates, err := security.DownloadCertificates(ctx, httpClient)
	if err != nil {
		return err
	}

	for _, cert := range certificates {
		certHash := certsdb.CertificateHash{
			Usage:     cert.Usage,
			Gateway:   gatewayConfig.Gateway,
			ValidFrom: cert.ValidFrom,
			ValidTo:   cert.ValidTo,
		}

		var derBytes []byte
		if derBytes, err = base64.StdEncoding.DecodeString(cert.Certificate); err != nil {
			return err
		}

		created, err := certsDB.StoreMinistryCertificate(
			ctx,
			certsdb.MinistryCertificateData{
				Certificate: certsdb.Certificate{
					CertificateHash: certHash,
				},
				DERBytes: derBytes,
			},
		)
		if err != nil {
			return err
		}

		if created {
			logging.CertsDBLogger.Debug("certyfikat ministerstwa finansów dodany do bazy", "usage", certHash.UsageAsString())
		}
	}

	return certsDB.Close(ctx)
}
