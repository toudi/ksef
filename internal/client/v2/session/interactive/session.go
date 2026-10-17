package interactive

import (
	"context"
	"errors"
	"ksef/core/certsdb"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/logging"

	HTTP "ksef/internal/http"

	"github.com/spf13/viper"
)

type Session struct {
	httpClient *HTTP.Client
	certsDB    *certsdb.CertsDB
	vip        *viper.Viper
}

var (
	ErrObtainSessionTokenTimeout = errors.New("timeout waiting for session token")
	ErrProbablyUsedSend          = errors.New("upload command probably used previously")
)

func NewSession(httpClient *HTTP.Client, certsDB *certsdb.CertsDB, vip *viper.Viper) *Session {
	return &Session{
		httpClient: httpClient,
		certsDB:    certsDB,
		vip:        vip,
	}
}

func (s *Session) UploadInvoices(
	ctx context.Context,
	payload types.UploadPayload,
) ([]*types.UploadSessionResult, error) {
	var result []*types.UploadSessionResult
	// v2 specs forces us to group invoices by their form code
	// on the other hand, it no longer forces us to send invoices through a 3rd party server
	for formCode, files := range payload {
		logging.InteractiveLogger.With("formCode", formCode).Info("przesyłanie faktur")
		uploadResult, err := s.uploadInvoicesForForm(ctx, formCode, files)
		if err != nil {
			return result, err
		}
		result = append(result, uploadResult)
	}

	return result, nil
}
