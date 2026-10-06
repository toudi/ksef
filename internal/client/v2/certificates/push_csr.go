package certificates

import (
	"context"
	baseHTTP "net/http"

	"ksef/internal/certsdb"
	"ksef/internal/http"
)

const (
	endpointEnrollments = "/api/v2/certificates/enrollments"
)

type enrollmentRequest struct {
	CertificateName string `json:"certificateName"`
	CertificateType string `json:"certificateType"`
	CSR             string `json:"csr"`
}

type enrollmentResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
}

func (m *Manager) PushCSR(ctx context.Context, cert *certsdb.Certificate) (*enrollmentResponse, error) {
	var resp enrollmentResponse
	eReq := enrollmentRequest{
		CertificateName: cert.Usage[0].Description(),
		CertificateType: string(cert.Usage[0]),
		CSR:             cert.CSRData,
	}

	_, err := m.httpClient.Request(
		ctx,
		http.RequestConfig{
			Method:          baseHTTP.MethodPost,
			Body:            eReq,
			ContentType:     http.JSON,
			Dest:            &resp,
			DestContentType: http.JSON,
			ExpectedStatus:  baseHTTP.StatusAccepted,
		},
		endpointEnrollments,
	)

	return &resp, err
}
