package interactive

import (
	"fmt"
	"time"
)

const sessionStatusEndpoint = "online/Session/Status/%s"

type sessionStatusType struct {
	ProcessingCode int `json:"processingCode"`
}

const VerifySecuritySuccess = 315

func (i *InteractiveSession) WaitForStatusCode(expectedCode int, maxRetries int) error {
	var sessionStatus sessionStatusType

	for j := 0; j < maxRetries; j += 1 {
		_, err := i.session.JSONRequest("GET", fmt.Sprintf(sessionStatusEndpoint, i.referenceNo), nil, &sessionStatus)
		if err != nil {
			return fmt.Errorf("error sending JSONRequest: %v", err)
		}
		if sessionStatus.ProcessingCode == expectedCode {
			return nil
		}
		fmt.Printf("session status code = %d\n", sessionStatus.ProcessingCode)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("maximum number of iterations reached")
}
