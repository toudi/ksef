package smart

import (
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/invoices/downloader/export"
	downloaderinterface "ksef/internal/client/v2/invoices/downloader/interface"
	"ksef/internal/client/v2/invoices/downloader/simple"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"log/slog"

	"github.com/spf13/viper"
)

type smartDownloader struct {
	vip        *viper.Viper
	certsDB    *certsdb.CertificatesDB
	httpClient *http.Client
	registry   *monthlyregistry.Registry
	params     invoices.DownloadParams
	logger     *slog.Logger
	// downloader instances (for reusing)
	simpleDownloader *simple.SimpleDownloader
	exportDownloader *export.ExportDownloader
}

func NewDownloader(
	vip *viper.Viper,
	certsDB *certsdb.CertificatesDB,
	httpClient *http.Client,
	registry *monthlyregistry.Registry,
	params invoices.DownloadParams,
	logger *slog.Logger,
) downloaderinterface.InvoiceDownloader {
	return &smartDownloader{
		vip:        vip,
		certsDB:    certsDB,
		httpClient: httpClient,
		registry:   registry,
		params:     params,
		logger:     logger,
	}
}

func (sd *smartDownloader) Close() error {
	if sd.exportDownloader != nil {
		return sd.exportDownloader.Close()
	}
	return nil
}
