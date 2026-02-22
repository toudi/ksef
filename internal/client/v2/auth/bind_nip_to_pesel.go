package auth

import (
	"context"
	"ksef/internal/http"
	baseHTTP "net/http"
)

const (
	endpointCreatePerson = "/v2/testdata/person"
)

type BindNipToPeselRequest struct {
	NIP         string `json:"nip"`
	PESEL       string `json:"pesel"`
	IsBailiff   bool   `json:"isBailiff"`
	Description string `json:"description"`
}

func BindNIPToPESEL(ctx context.Context, httpClient *http.Client, nip, pesel string) error {
	_, err := httpClient.Request(
		ctx,
		http.RequestConfig{
			Body: BindNipToPeselRequest{
				NIP:         nip,
				PESEL:       pesel,
				IsBailiff:   false,
				Description: "JDG: Powiązanie NIP " + nip + " z numerem PESEL " + pesel,
			},
			ContentType:    http.JSON,
			Method:         baseHTTP.MethodPost,
			ExpectedStatus: baseHTTP.StatusOK,
		},
		endpointCreatePerson,
	)

	return err
}
