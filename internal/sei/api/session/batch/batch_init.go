package batch

import "ksef/internal/sei/api/client"

type BatchSession struct {
	apiClient *client.APIClient
}

func BatchSessionInit(client *client.APIClient) *BatchSession {
	return &BatchSession{apiClient: client}
}
