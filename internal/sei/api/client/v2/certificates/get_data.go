package certificates

import (
	"context"
	"ksef/internal/http"
	baseHTTP "net/http"
)

const (
	endpointEnrollmentsData = "/api/v2/certificates/enrollments/data"
)

type EnrollmentsData struct {
	CommonName             string  `json:"commonName"`
	CountryName            string  `json:"countryName"`
	GivenName              *string `json:"givenName"`
	Surname                *string `json:"surname"`
	SerialNumber           *string `json:"serialNumber"`
	UniqueIdentifier       *string `json:"uniqueIdentifier"`
	OrganizationName       *string `json:"organizationName"`
	OrganizationIdentifier *string `json:"organizationIdentifier"`
}

func (m *Manager) GetEnrollmentsData(ctx context.Context) (*EnrollmentsData, error) {
	var ed EnrollmentsData

	_, err := m.httpClient.Request(
		ctx,
		http.RequestConfig{
			Method:          baseHTTP.MethodGet,
			Dest:            &ed,
			DestContentType: http.JSON,
		},
		endpointEnrollmentsData,
	)

	if err != nil {
		return nil, err
	}

	return &ed, nil
}
