package status

import "ksef/internal/http"

type SessionStatusChecker struct {
	httpClient *http.Client
}

func NewSessionStatusChecker(httpClient *http.Client) *SessionStatusChecker {
	return &SessionStatusChecker{
		httpClient: httpClient,
	}
}
