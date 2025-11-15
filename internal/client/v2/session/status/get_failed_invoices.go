package status

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"net/http"
)

const endpointFailedInvoices = "/api/v2/sessions/%s/invoices/failed"

type FailedInvoiceDetails struct {
	OrdinalNumber int
	Details       []string
}

type FailedInvoiceStatus struct {
	Description string   `json:"description"`
	Details     []string `json:"details"`
}
type FailedInvoiceInfo struct {
	OrdinalNumber int                 `json:"ordinalNumber"`
	Status        FailedInvoiceStatus `json:"status"`
}

type FailedInvoicesResponse struct {
	ContinuationToken string              `json:"continuationToken"`
	Invoices          []FailedInvoiceInfo `json:"invoices"`
}

func GetFailedInvoiceList(
	ctx context.Context,
	httpClient *HTTP.Client,
	uploadSessionId string,
) ([]FailedInvoiceDetails, error) {
	var failedInvoices []FailedInvoiceDetails
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
			var failedInvoiceDetails = FailedInvoiceDetails{
				// important: OrdinalNumbers are 1-based
				OrdinalNumber: invoiceInfo.OrdinalNumber - 1,
				Details:       []string{invoiceInfo.Status.Description},
			}
			failedInvoiceDetails.Details = append(failedInvoiceDetails.Details, invoiceInfo.Status.Details...)
			failedInvoices = append(failedInvoices, failedInvoiceDetails)
		}
	}

	return failedInvoices, nil
}
