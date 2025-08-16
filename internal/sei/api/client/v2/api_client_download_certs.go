package v2

import (
	"context"
	"net/http"
	"time"
)

const downloadCertsEndpoint = "/security/public-key-certificates"

type CertificateRow struct {
	CertificateBase64 string    `json:"certificate"`
	ValidTo           time.Time `json:"validTo"`
	Usage             []string  `json:"usage"`
}

func (a *APIClient) DownloadCerts() error {
	ctx, cancel := context.WithTimeout(a.ctx, 15*time.Second)
	defer cancel()

	req, err := a.NewRequest(ctx, http.MethodGet, downloadCertsEndpoint, nil)
	if err != nil {
		return err
	}

	var certificates []CertificateRow
	err = a.DoJSONResponse(req, &certificates)
	if err != nil {
		return err
	}
	return nil
}
