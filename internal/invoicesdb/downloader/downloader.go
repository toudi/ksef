package downloader

import (
	v2 "ksef/internal/client/v2"

	"github.com/spf13/viper"
)

type Downloader struct{}

func NewDownloader(vip *viper.Viper, client *v2.APIClient) *Downloader {
	return &Downloader{}
}
