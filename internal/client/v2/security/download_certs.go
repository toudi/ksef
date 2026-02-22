package security

import (
	"context"
	"ksef/internal/certsdb"
	"ksef/internal/http"
	"time"
)

const endpointDownloadCertificates = "/v2/security/public-key-certificates"

type certificateRow struct {
	Certificate string          `json:"certificate"`
	ValidFrom   time.Time       `json:"validFrom"`
	ValidTo     time.Time       `json:"validTo"`
	Usage       []certsdb.Usage `json:"usage"`
}

func DownloadCertificates(ctx context.Context, client *http.Client) ([]certificateRow, error) {
	var certificates []certificateRow

	_, err := client.Request(
		ctx,
		http.RequestConfig{
			Dest:            &certificates,
			DestContentType: http.JSON,
		},
		endpointDownloadCertificates,
	)

	return certificates, err
}
