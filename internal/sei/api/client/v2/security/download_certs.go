package security

import (
	"context"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/http"
)

const endpointDownloadCertificates = "/api/v2/security/public-key-certificates"

type certificateRow struct {
	Certificate string          `json:"certificate"`
	Usage       []certsdb.Usage `json:"usage"`
}

func DownloadCertificates(ctx context.Context, client *http.Client, cfg config.APIConfig) error {
	var certificates []certificateRow

	_, err := client.Request(
		ctx,
		http.RequestConfig{
			Dest:            &certificates,
			DestContentType: http.JSON,
		},
		endpointDownloadCertificates,
	)

	if err != nil {
		return err
	}

	for _, certificate := range certificates {
		if err := cfg.CertificatesDB.AddCertificate(certificate.Certificate, cfg.Environment.Environment, certificate.Usage); err != nil {
			return err
		}
	}

	return nil
}
