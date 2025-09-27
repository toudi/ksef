package batch

import (
	"ksef/internal/config"
	HTTP "ksef/internal/http"
	"ksef/internal/registry"
)

type Session struct {
	registry   *registry.InvoiceRegistry
	httpClient *HTTP.Client
	apiConfig  config.APIConfig
}

func NewSession(
	httpClient *HTTP.Client,
	registry *registry.InvoiceRegistry,
	apiConfig config.APIConfig,
) *Session {
	return &Session{
		httpClient: httpClient,
		registry:   registry,
		apiConfig:  apiConfig,
	}
}
