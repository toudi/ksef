package interactive

import (
	"fmt"
)

const endpointGeneratePaymentId = "online/Payment/Identifier/Request"

type paymentIdRequest struct {
	SEIInvoiceRefNos []string `json:"ksefReferenceNumberList"`
}

type paymentIdResponse struct {
	RefNo string `json:"paymentIdentifier"`
}

func (i *InteractiveSession) GeneratePaymentId(seiInvoiceRefNos []string) (string, error) {
	if err := i.WaitForStatusCode(VerifySecuritySuccess, 15); err != nil {
		return "", fmt.Errorf("unable to obtain successful status code within 15 seconds: %v", err)
	}
	var request paymentIdRequest
	request.SEIInvoiceRefNos = seiInvoiceRefNos
	var response paymentIdResponse

	httpResponse, err := i.session.JSONRequest("POST", endpointGeneratePaymentId, request, &response)
	if httpResponse.StatusCode != 201 {
		return "", fmt.Errorf("unexpected response from generating payment ID: %d vs 201", httpResponse.StatusCode)
	}
	if err != nil {
		return "", fmt.Errorf("unable to generate payment UID: %v", err)
	}

	return response.RefNo, nil
}
