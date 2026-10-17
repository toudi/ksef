package export

import (
	"ksef/core/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

type ExportDownloader struct {
	vip            *viper.Viper
	httpClient     *http.Client
	registry       *monthlyregistry.Registry
	params         invoices.DownloadParams
	certsDB        *certsdb.CertsDB
	archiveHandler *archiveHandler
}

func NewDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertsDB,
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
) *ExportDownloader {
	return &ExportDownloader{
		vip:        vip,
		certsDB:    certsDB,
		httpClient: httpClient,
		registry:   registry,
		params:     params,
	}
}

func (ed *ExportDownloader) Close() error {
	if ed.archiveHandler != nil {
		return ed.archiveHandler.Close()
	}
	return nil
}
