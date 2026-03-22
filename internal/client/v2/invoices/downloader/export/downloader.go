package export

import (
	"ksef/internal/certsdb"
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
	certsDB        *certsdb.CertificatesDB
	archiveHandler *archiveHandler
}

func NewDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertificatesDB,
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
