package status

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"ksef/internal/utils"
	"net/http"
)

const endpointSessionInvoices = "/api/v2/sessions/%s/invoices"

type InvoiceStatus struct {
	Code        int      `json:"code"`
	Description string   `json:"description"`
	Details     []string `json:"details"`
}

func (is InvoiceStatus) Successful() bool {
	return is.Code == 200
}

type InvoiceInfo struct {
	OrdinalNumber  int           `json:"ordinalNumber"`
	ChecksumBase64 string        `json:"invoiceHash"` // checksum as returned by KSeF gateway, i.e. raw bytes of sha256sum encoded with base64
	Checksum       string        // checksum as a simple hex string, i.e. result of sha256sum <file>
	InvoiceNumber  string        `json:"invoiceNumber"`
	KSeFRefNo      string        `json:"ksefNumber"`
	Status         InvoiceStatus `json:"status"`
}

type SessionInvoicesResponse struct {
	ContinuationToken string        `json:"continuationToken"`
	Invoices          []InvoiceInfo `json:"invoices"`
}

func (sc *SessionStatusChecker) GetInvoiceList(
	ctx context.Context,
	uploadSessionId string,
) ([]InvoiceInfo, error) {
	var sessionInvoices []InvoiceInfo
	var err error

	var invoicesResponse SessionInvoicesResponse

	var finished bool = false
	var continuationToken string

	for !finished {
		var headers = map[string]string{}
		if continuationToken != "" {
			headers["x-continuation-token"] = continuationToken
		}

		_, err = sc.httpClient.Request(
			ctx,
			HTTP.RequestConfig{
				Headers:         headers,
				ContentType:     HTTP.JSON,
				Dest:            &invoicesResponse,
				DestContentType: HTTP.JSON,
				ExpectedStatus:  http.StatusOK,
				Method:          http.MethodGet,
			},
			fmt.Sprintf(endpointSessionInvoices, uploadSessionId),
		)

		if err != nil {
			return nil, err
		}

		finished = invoicesResponse.ContinuationToken == ""

		for _, invoiceInfo := range invoicesResponse.Invoices {
			var allDetails []string
			if !invoiceInfo.Status.Successful() {
				allDetails = append(allDetails, invoiceInfo.Status.Description)
				allDetails = append(allDetails, invoiceInfo.Status.Details...)
				invoiceInfo.Status.Details = allDetails
			}

			// unfortunetely, KSeF API returns checksums in base64 form which is space-efficient,
			// but makes it painful to compare with local checksums generated with command-line
			// utilities, therefore internally we're keeping the hex-encoded checksum
			invoiceInfo.Checksum = utils.Base64ToHex(invoiceInfo.ChecksumBase64)
			sessionInvoices = append(sessionInvoices, invoiceInfo)
		}
	}

	return sessionInvoices, nil
}
