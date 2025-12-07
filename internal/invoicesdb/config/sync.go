package config

import (
	"ksef/internal/client/v2/types/invoices"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type SyncConfig struct {
	Uploader   uploaderconfig.UploaderConfig
	Downloader invoices.DownloadParams
}

const (
	cfgKeyWaitForStatus = "upload.wait"
	cfgKeyWaitTimeout   = "upload.wait.timeout"
	cfgKeyDownloadUpo   = "upo"
	cfgKeyConvertUPOPdf = "upo.pdf"
)

func SyncFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyWaitForStatus, false, "czekaj na zakończenie wysyłki")
	flagSet.Duration(cfgKeyWaitTimeout, time.Duration(0), "maksymalny czas oczekiwania na rezultat wysyłki")
	flagSet.Bool(cfgKeyDownloadUpo, false, "pobierz UPO po przetworzeniu sesji")
	flagSet.Bool(cfgKeyConvertUPOPdf, false, "konwertuj UPO do PDF")
	downloaderconfig.DownloaderFlags(flagSet, "download")

	flagSet.SortFlags = false
}

func GetSyncConfig(vip *viper.Viper) (SyncConfig, error) {
	if downloaderConfig, err := downloaderconfig.GetDownloaderConfig(vip); err != nil {
		return SyncConfig{}, err
	} else {
		return SyncConfig{
			Uploader: uploaderconfig.UploaderConfig{
				WaitForStatus: vip.GetBool(cfgKeyWaitForStatus),
				WaitTimeout:   vip.GetDuration(cfgKeyWaitTimeout),
				DownloadUPO:   vip.GetBool(cfgKeyDownloadUpo),
				SaveUPOAsPDF:  vip.GetBool(cfgKeyConvertUPOPdf),
			},
			Downloader: downloaderConfig,
		}, nil
	}
}
