package resources

import (
	"ksef/core/downloader"

	"github.com/spf13/viper"
)

var instance *downloader.Config

func DownloaderConfig(vip *viper.Viper) (downloader.Config, error) {
	if instance != nil {
		return *instance, nil
	}

	return downloader.Config{}, nil
}
