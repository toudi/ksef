package config

import (
	"ksef/internal/client/v2/types/invoices"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type SyncConfig struct {
	Uploader   uploaderconfig.UploaderConfig
	Downloader invoices.DownloadParams
}

func SyncFlags(flagSet *pflag.FlagSet) {
	UploaderFlags(flagSet)
	downloaderconfig.DownloaderFlags(flagSet, "download")

	flagSet.SortFlags = false
}

func GetSyncConfig(vip *viper.Viper) (SyncConfig, error) {
	if downloaderConfig, err := downloaderconfig.GetDownloaderConfig(vip); err != nil {
		return SyncConfig{}, err
	} else {
		return SyncConfig{
			Uploader:   GetUploaderConfig(vip),
			Downloader: downloaderConfig,
		}, nil
	}
}
