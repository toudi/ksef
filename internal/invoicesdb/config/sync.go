package config

import (
	"ksef/internal/client/v2/types/invoices"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	downloaderFlagsPrefix = "download"
)

type SyncConfig struct {
	Uploader            uploaderconfig.UploaderConfig
	Downloader          invoices.DownloadParams
	StatusCheckerConfig statuscheckerconfig.StatusCheckerConfig
}

func SyncFlags(flagSet *pflag.FlagSet) {
	uploaderconfig.UploaderFlags(flagSet)
	statuscheckerconfig.StatusCheckerFlags(flagSet)
	downloaderconfig.DownloaderFlags(flagSet, downloaderFlagsPrefix)

	flagSet.SortFlags = false
}

func GetSyncConfig(vip *viper.Viper) (SyncConfig, error) {
	if downloaderConfig, err := downloaderconfig.GetDownloaderConfig(vip, downloaderFlagsPrefix); err != nil {
		return SyncConfig{}, err
	} else {
		return SyncConfig{
			Uploader:            uploaderconfig.GetUploaderConfig(vip),
			Downloader:          downloaderConfig,
			StatusCheckerConfig: statuscheckerconfig.GetStatusCheckerConfig(vip),
		}, nil
	}
}
