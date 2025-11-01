package interactive

import (
	"context"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/registry"
)

type Session struct {
	registry   *registry.InvoiceRegistry
	httpClient *HTTP.Client
	apiConfig  config.APIConfig
	certsDB    *certsdb.CertificatesDB
}

var (
	ErrObtainSessionTokenTimeout = errors.New("timeout waiting for session token")
	ErrProbablyUsedSend          = errors.New("upload command probably used previously")
)

func NewSession(httpClient *HTTP.Client, registry *registry.InvoiceRegistry, certsDB *certsdb.CertificatesDB) *Session {
	return &Session{
		httpClient: httpClient,
		registry:   registry,
		certsDB:    certsDB,
	}
}

func (s *Session) UploadInvoices(ctx context.Context, params UploadParams) error {
	// at this point the collection has already been initialized and retrieved so no need for checking the error
	collection, _ := s.registry.InvoiceCollection()

	// v2 specs forces us to group invoices by their form code
	// on the other hand, it no longer forces us to send invoices through a 3rd party server
	for formCode, files := range collection.Files {
		logging.InteractiveLogger.With("formCode", formCode).Info("przesyłanie faktur")
		if err := s.uploadInvoicesForForm(ctx, formCode, files); err != nil {
			return err
		}
	}

	return nil
}
