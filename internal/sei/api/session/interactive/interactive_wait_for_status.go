package interactive

import (
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client"
	"time"
)

const sessionStatusEndpoint = "/api/online/Session/Status/%s"

type sessionStatusType struct {
	ProcessingCode int `json:"processingCode"`
}

const VerifySecuritySuccess = 315

func (i *InteractiveSession) WaitForStatusCode(expectedCode int, maxRetries int) error {
	var sessionStatus sessionStatusType

	for j := 0; j < maxRetries; j += 1 {
		_, err := i.session.JSONRequest(
			client.JSONRequestParams{
				Method:   "GET",
				Endpoint: fmt.Sprintf(sessionStatusEndpoint, i.referenceNo),
				Payload:  nil,
				Response: &sessionStatus,
				Logger:   logging.InteractiveHTTPLogger,
			},
		)
		if err != nil {
			return fmt.Errorf("error sending JSONRequest: %v", err)
		}
		if sessionStatus.ProcessingCode == expectedCode {
			return nil
		}

		logging.InteractiveLogger.Debug(
			"InteractiveSession::WaitForStatusCode",
			"status",
			sessionStatus.ProcessingCode,
		)

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("maximum number of iterations reached")
}
