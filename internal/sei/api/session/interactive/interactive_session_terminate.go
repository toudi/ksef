package interactive

import (
	"fmt"
	"io"
	"ksef/internal/logging"
	"net/http"
)

const endpointLogout = "/api/online/Session/Terminate"

func (i *InteractiveSession) Logout() error {
	terminateRequest, err := i.session.Request(
		"GET",
		endpointLogout,
		nil,
		logging.InteractiveHTTPLogger,
	)
	if err != nil {
		return fmt.Errorf("unable to perform request: %v", err)
	}

	terminateResponse, err := http.DefaultClient.Do(terminateRequest)
	if err != nil || terminateResponse.StatusCode/100 != 2 {
		defer terminateResponse.Body.Close()
		terminateResponseBody, _ := io.ReadAll(terminateResponse.Body)
		return fmt.Errorf(
			"error finishing session: statuscode=%d, err=%v\n%s",
			terminateResponse.StatusCode,
			err,
			string(terminateResponseBody),
		)
	}

	return nil
}
