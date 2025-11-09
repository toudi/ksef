package certificates

import (
	"context"
	"ksef/internal/http"
	baseHTTP "net/http"
)

const (
	endpointCertificateRetrieve = "/api/v2/certificates/retrieve"
)

type CertificateDetails struct {
	Certificate string `json:"certificate"`
}

type DownloadCertificatesResponse struct {
	Certificates []CertificateDetails `json:"certificates"`
}

type DownloadCertificatesRequest struct {
	SerialNumbers []string `json:"certificateSerialNumbers"`
}

func (m *Manager) DownloadCertificate(
	ctx context.Context,
	serialNumber string,
) (DownloadCertificatesResponse, error) {
	var resp DownloadCertificatesResponse
	var req = DownloadCertificatesRequest{
		SerialNumbers: []string{serialNumber},
	}

	_, err := m.httpClient.Request(
		ctx,
		http.RequestConfig{
			Method:          baseHTTP.MethodPost,
			ContentType:     http.JSON,
			Body:            req,
			Dest:            &resp,
			DestContentType: http.JSON,
			ExpectedStatus:  baseHTTP.StatusOK,
		},
		endpointCertificateRetrieve,
	)

	return resp, err
}
