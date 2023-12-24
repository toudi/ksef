package interactive

import (
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client"
)

const endpointGeneratePaymentId = "/api/online/Payment/Identifier/Request"

type paymentIdRequest struct {
	SEIInvoiceRefNos []string `json:"ksefReferenceNumberList"`
}

type paymentIdResponse struct {
	RefNo string `json:"paymentIdentifier"`
}

func (i *InteractiveSession) GeneratePaymentId(seiInvoiceRefNos []string) (string, error) {
	var request paymentIdRequest
	request.SEIInvoiceRefNos = seiInvoiceRefNos
	var response paymentIdResponse

	httpResponse, err := i.session.JSONRequest(
		client.JSONRequestParams{
			Method:   "POST",
			Endpoint: endpointGeneratePaymentId,
			Payload:  request,
			Response: &response,
			Logger:   logging.InteractiveHTTPLogger,
		},
	)
	if httpResponse.StatusCode != 201 {
		return "", fmt.Errorf(
			"unexpected response from generating payment ID: %d vs 201",
			httpResponse.StatusCode,
		)
	}
	if err != nil {
		return "", fmt.Errorf("unable to generate payment UID: %v", err)
	}

	return response.RefNo, nil
}
