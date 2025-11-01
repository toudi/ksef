package batch

import (
	"ksef/internal/certsdb"
	HTTP "ksef/internal/http"
	"ksef/internal/registry"
)

type Session struct {
	registry   *registry.InvoiceRegistry
	httpClient *HTTP.Client
	certsDB    *certsdb.CertificatesDB
}

func NewSession(
	httpClient *HTTP.Client,
	registry *registry.InvoiceRegistry,
	certsDB *certsdb.CertificatesDB,
) *Session {
	return &Session{
		httpClient: httpClient,
		registry:   registry,
		certsDB:    certsDB,
	}
}
