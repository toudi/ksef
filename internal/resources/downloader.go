package resources

import (
	"ksef/core/downloader"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"

	"github.com/spf13/viper"
)

var instance *downloader.Config

func DownloaderConfig(vip *viper.Viper) (downloader.Config, error) {
	if instance != nil {
		return *instance, nil
	}

	params, err := downloaderconfig.GetDownloaderConfig(vip, "")
	if err != nil {
		return downloader.Config{}, err
	}
	instance = &params

	return *instance, nil
}
