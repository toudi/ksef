package smart

import (
	"ksef/core/certsdb"
	"ksef/core/downloader"
	"ksef/internal/client/v2/invoices/downloader/export"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/invoices/downloader/simple"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

type smartDownloader struct {
	vip        *viper.Viper
	certsDB    *certsdb.CertsDB
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     downloader.Config
	// downloader instances (for reusing)
	simpleDownloader *simple.SimpleDownloader
	exportDownloader *export.ExportDownloader
}

func NewDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertsDB,
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
) downloaderinterface.InvoiceDownloader {
	return &smartDownloader{
		vip:        vip,
		certsDB:    certsDB,
		httpClient: httpClient,
		registry:   registry,
	}
}

func (sd *smartDownloader) Close() error {
	if sd.exportDownloader != nil {
		return sd.exportDownloader.Close()
	}
	return nil
}
