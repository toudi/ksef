package status

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"net/http"
)

const endpointFailedInvoices = "/api/v2/sessions/{referenceNumber}/invoices/failed"

type FailedInvoiceInfo struct {
	OrdinalNumber int `json:"ordinalNumber"`
}

type FailedInvoicesResponse struct {
	ContinuationToken string              `json:"continuationToken"`
	Invoices          []FailedInvoiceInfo `json:"invoices"`
}

func GetFailedInvoiceList(ctx context.Context, httpClient HTTP.Client, uploadSessionId string) ([]int, error) {
	var failedInvoices []int
	var err error

	var failedInvoicesResponse FailedInvoicesResponse

	var finished bool = false
	var continuationToken string

	for !finished {
		var headers = map[string]string{}
		if continuationToken != "" {
			headers["x-continuation-token"] = continuationToken
		}

		_, err = httpClient.Request(
			ctx,
			HTTP.RequestConfig{
				Headers:         headers,
				ContentType:     HTTP.JSON,
				Dest:            &failedInvoicesResponse,
				DestContentType: HTTP.JSON,
				ExpectedStatus:  http.StatusOK,
				Method:          http.MethodGet,
			},
			fmt.Sprintf(endpointFailedInvoices, uploadSessionId),
		)

		if err != nil {
			return nil, err
		}

		finished = failedInvoicesResponse.ContinuationToken == ""

		for _, invoiceInfo := range failedInvoicesResponse.Invoices {
			failedInvoices = append(failedInvoices, invoiceInfo.OrdinalNumber)
		}
	}

	return failedInvoices, nil
}
