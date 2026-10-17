package certificates

import (
	"context"
	"ksef/core/certsdb"
	"ksef/internal/http"
	baseHTTP "net/http"
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

func (m *Manager) PushCSR(ctx context.Context, cert certsdb.Certificate) (*enrollmentResponse, error) {
	var resp enrollmentResponse

	// TODO: maybe we should check if the length of Usage is exactly 1 ?
	// on the other hand, it could be checked during creation of the certificate.
	// check that after refactor is done.
	certUsageDescr, err := cert.Usage[0].Description()
	if err != nil {
		return nil, err
	}

	eReq := enrollmentRequest{
		CertificateName: certUsageDescr,
		CertificateType: string(cert.Usage[0]),
		CSR:             cert.CSRData,
	}

	_, err = m.httpClient.Request(
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
