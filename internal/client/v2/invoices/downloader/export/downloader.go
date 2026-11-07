package export

import (
	"ksef/core/certsdb"
	"ksef/core/downloader"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

type ExportDownloader struct {
	vip            *viper.Viper
	httpClient     *http.Client
	registry       *monthlyregistry.Registry
	params         downloader.Config
	certsDB        *certsdb.CertsDB
	archiveHandler *archiveHandler
}

func NewDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertsDB,
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
) *ExportDownloader {
	return &ExportDownloader{
		vip:        vip,
		certsDB:    certsDB,
		httpClient: httpClient,
		registry:   registry,
	}
}

func (ed *ExportDownloader) Close() error {
	if ed.archiveHandler != nil {
		return ed.archiveHandler.Close()
	}
	return nil
}
