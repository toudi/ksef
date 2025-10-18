package certificates

import (
	"context"
	"errors"
	"fmt"
	"ksef/internal/certsdb"
	"ksef/internal/http"
	baseHTTP "net/http"
)

const (
	endpointEnrollmentStatus     = "/api/v2/certificates/enrollments/%s"
	enrollmentSuccess        int = 200
)

var (
	errInvalidEnrollment = errors.New("nieprawidłowy wniosek (brak numeru referencyjnego)")
)

type status struct {
	Code        int       `json:"code"`
	Description string    `json:"description"`
	Details     *[]string `json:"details"`
}
type enrollmentStatusResponse struct {
	Status       status
	SerialNumber *string `json:"certificateSerialNumber"`
}

func (m *Manager) GetEnrollmentStatus(ctx context.Context, cert *certsdb.Certificate) (*enrollmentStatusResponse, error) {
	if cert.ReferenceNumber == "" {
		return nil, errInvalidEnrollment
	}

	var resp enrollmentStatusResponse

	_, err := m.httpClient.Request(
		ctx,
		http.RequestConfig{
			Method:          baseHTTP.MethodGet,
			Dest:            &resp,
			DestContentType: http.JSON,
			ExpectedStatus:  baseHTTP.StatusOK,
		},
		fmt.Sprintf(endpointEnrollmentStatus, cert.ReferenceNumber),
	)

	return &resp, err
}
