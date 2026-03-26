package export

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"log/slog"

	"github.com/spf13/viper"
)

type ExportDownloader struct {
	vip            *viper.Viper
	httpClient     *http.Client
	registry       *monthlyregistry.Registry
	params         invoices.DownloadParams
	certsDB        *certsdb.CertificatesDB
	archiveHandler *archiveHandler
	logger         *slog.Logger
}

func NewDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertificatesDB,
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
	logger *slog.Logger,
) *ExportDownloader {
	return &ExportDownloader{
		vip:        vip,
		certsDB:    certsDB,
		httpClient: httpClient,
		registry:   registry,
		params:     params,
		logger:     logger,
	}
}

func (ed *ExportDownloader) Close() error {
	if ed.archiveHandler != nil {
		return ed.archiveHandler.Close()
	}
	return nil
}
